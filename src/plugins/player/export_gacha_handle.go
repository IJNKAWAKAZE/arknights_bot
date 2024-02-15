package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/xuri/excelize/v2"
	"io"
	"os"
	"time"
)

func ExportGachaHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

	var userAccount account.UserAccount
	var players []account.UserPlayer

	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, "未查询到绑定账号，请先进行绑定。")
		bot.Arknights.Send(sendMessage)
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
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s,%d,%s,%d", "player", OP_EXPORT, userId, player.Uid, messageId)),
			))
		}
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要导出的角色")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return true, nil
	} else {
		// 绑定单个角色
		return Export(players[0].Uid, userAccount, chatId)
	}
	return true, nil
}

func Export(uid string, account account.UserAccount, chatId int64) (bool, error) {
	var userGacha []UserGacha
	res := utils.GetUserGacha(account.UserNumber, uid).Scan(&userGacha)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "不存在抽卡记录！")
		bot.Arknights.Send(sendMessage)
		return true, nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_document")
	bot.Arknights.Send(sendAction)

	f := excelize.NewFile()
	// 设置单元格的值
	f.SetSheetRow("Sheet1", "A1", &[]interface{}{"卡池", "干员", "星级", "New", "时间"})
	f.SetColWidth("Sheet1", "A", "E", 18)
	for i, gacha := range userGacha {
		f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", i+2), &[]interface{}{gacha.PoolName, gacha.CharName, gacha.Rarity + 1, gacha.IsNew, time.Unix(gacha.Ts, 0).Format("2006-01-02 15:04:05")})
	}
	fileName := time.Now().Format("2006-01-02") + ".xlsx"
	// 根据指定路径保存文件
	if err := f.SaveAs(fileName); err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成文件失败！")
		bot.Arknights.Send(sendMessage)
	}

	file, _ := os.Open(fileName)
	b, _ := io.ReadAll(file)
	file.Close()
	os.Remove(fileName)
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: b, Name: fileName})
	sendDocument.Caption = "抽卡记录导出成功！"
	bot.Arknights.Send(sendDocument)
	return true, nil
}
