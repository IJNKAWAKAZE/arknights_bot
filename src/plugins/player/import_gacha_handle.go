package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"arknights_bot/utils/telebot"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"github.com/tidwall/sjson"
	"io"
	"strconv"
	"strings"
)

var ImportFile = make(map[int64]string)

type ImportGachaData struct {
	Data map[string]struct {
		C [][]interface{} `json:"c"`
		P string          `json:"p"`
	} `json:"data"`
}

func ImportGachaHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID

	var userAccount account.UserAccount

	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, "未查询到绑定账号，请先进行绑定。")
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	sendMessage := tgbotapi.NewMessage(chatId, "请将[网站](https://arkgacha.kwer.top/)导出的抽卡数据粘贴到txt文本中发送或使用 /cancel 指令取消操作。")
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	bot.Arknights.Send(sendMessage)
	telebot.WaitMessage[chatId] = "importGacha"
	return true, nil
}

// ImportGacha 导入抽卡记录
func ImportGacha(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	doc := update.Message.Document
	if doc != nil && strings.HasSuffix(doc.FileName, ".txt") {
		delete(telebot.WaitMessage, chatId)
		ImportFile[userId] = doc.FileID
		var userAccount account.UserAccount
		var players []account.UserPlayer

		res := utils.GetAccountByUserId(userId).Scan(&userAccount)
		if res.RowsAffected == 0 {
			// 未绑定账号
			sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("未查询到绑定账号，请先进行[绑定](https://t.me/%s)。", viper.GetString("bot.name")))
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			sendMessage.ReplyToMessageID = messageId
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(chatId, messageId, 5)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		}

		// 获取绑定角色
		res = utils.GetPlayersByUserId(userId).Scan(&players)
		if res.RowsAffected == 0 {
			sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(chatId, messageId, 5)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		}

		if len(players) > 1 {
			// 绑定多个角色进行选择
			var buttons [][]tgbotapi.InlineKeyboardButton
			for _, player := range players {
				buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s,%d,%s,%d", "player", OP_IMPORT, userId, player.Uid, messageId)),
				))
			}
			inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
				buttons...,
			)
			sendMessage := tgbotapi.NewMessage(chatId, "请选择要导入的角色")
			sendMessage.ReplyMarkup = inlineKeyboardMarkup
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		} else {
			// 绑定单个角色
			return Import(players[0].Uid, userAccount, chatId, doc.FileID, utils.GetFullName(update.Message.From))
		}
	}

	sendMessage := tgbotapi.NewMessage(chatId, "导入文件格式错误！")
	sendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(sendMessage)
	return true, nil
}

func Import(uid string, account account.UserAccount, chatId int64, fileId string, name string) (bool, error) {
	var importGachaData ImportGachaData
	f, _ := utils.DownloadFile(fileId)
	data, _ := io.ReadAll(f)
	j, _ := sjson.SetRaw("{}", "data", string(data))
	err := json.Unmarshal([]byte(j), &importGachaData)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "解析抽卡记录失败！")
		bot.Arknights.Send(sendMessage)
		return true, err
	}

	go addGacha(importGachaData, account.UserNumber, uid, name)
	sendMessage := tgbotapi.NewMessage(chatId, "抽卡记录导入成功！")
	bot.Arknights.Send(sendMessage)
	return true, nil
}

func addGacha(importGachaData ImportGachaData, userNumber int64, uid string, name string) {
	// 遍历抽卡记录
	for k, d := range importGachaData.Data {
		key, _ := strconv.ParseInt(k, 10, 64)
		var gacha UserGacha
		res := bot.DBEngine.Raw("select * from user_gacha where user_number = ? and uid = ? and ts = ?", userNumber, uid, key).Scan(&gacha)
		if res.RowsAffected == 0 {
			// 同步抽卡数据
			for i, c := range d.C {
				id, _ := gonanoid.New(32)
				n := strconv.Itoa(int(c[2].(float64)))
				isNew, _ := strconv.ParseBool(n)
				userGacha := UserGacha{
					Id:         id,
					UserName:   name,
					UserNumber: userNumber,
					Uid:        uid,
					PoolName:   d.P,
					PoolOrder:  i + 1,
					CharName:   c[0].(string),
					IsNew:      isNew,
					Rarity:     int64(c[1].(float64)),
					Ts:         key,
				}
				bot.DBEngine.Table("user_gacha").Create(&userGacha)
			}
		}
	}
}
