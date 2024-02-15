package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"log"
	"time"
)

type UserGacha struct {
	Id         string    `json:"id" gorm:"primaryKey"`
	UserName   string    `json:"userName"`
	UserNumber int64     `json:"userNumber"`
	Uid        string    `json:"uid"`
	PoolName   string    `json:"poolName"`
	PoolOrder  int       `json:"poolOrder"`
	CharName   string    `json:"charName"`
	IsNew      bool      `json:"isNew"`
	Rarity     int64     `json:"rarity"`
	Ts         int64     `json:"ts"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark     string    `json:"remark"`
}

// SyncGachaHandle 同步抽卡记录
func SyncGachaHandle(players []account.UserPlayer, userAccount account.UserAccount, chatId int64, userId int64, messageId int) (bool, error) {
	if len(players) > 1 {
		// 绑定多个角色进行选择
		var buttons [][]tgbotapi.InlineKeyboardButton
		for _, player := range players {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s,%d,%s,%d", "player", OP_SYNC, userId, player.Uid, messageId)),
			))
		}
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要同步的角色")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	} else {
		// 绑定单个角色
		return SyncGacha(players[0].Uid, userAccount, chatId, messageId, userAccount.UserName)
	}
	return true, nil
}

func SyncGacha(uid string, userAccount account.UserAccount, chatId int64, messageId int, name string) (bool, error) {
	token := userAccount.HypergryphToken
	channelId := "1"
	var userPlayer account.UserPlayer
	utils.GetPlayerByUserId(userAccount.UserNumber, uid).Scan(&userPlayer)
	if userPlayer.ServerName == "b服" || userPlayer.ServerName == "bilibili服" {
		token = userAccount.BToken
		channelId = "2"
		// BToken为空设置BToken
		if token == "" {
			sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("B服Token未设置，请先进行[设置](https://t.me/%s)。", viper.GetString("bot.name")))
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			sendMessage.ReplyToMessageID = messageId
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(chatId, messageId, 5)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		}
	}
	// 获取角色抽卡记录
	chars, err := skland.GetPlayerGacha(token, channelId)
	if err != nil {
		log.Println(err)
		sendMessage := tgbotapi.NewMessage(chatId, "token可能已失效请重设token。")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, err
	}

	// 获取上次更新时间
	var lastUpdate int64
	bot.DBEngine.Raw("select ts from user_gacha where user_number = ? and uid = ? order by ts desc limit 1", userAccount.UserNumber, uid).Scan(&lastUpdate)

	// 同步抽卡数据
	for _, c := range chars {
		if c.Ts > lastUpdate {
			id, _ := gonanoid.New(32)
			userGacha := UserGacha{
				Id:         id,
				UserName:   name,
				UserNumber: userAccount.UserNumber,
				Uid:        uid,
				PoolName:   c.PoolName,
				PoolOrder:  c.PoolOrder,
				CharName:   c.Name,
				IsNew:      c.IsNew,
				Rarity:     c.Rarity,
				Ts:         c.Ts,
			}
			go bot.DBEngine.Table("user_gacha").Create(&userGacha)
		}
	}

	sendMessage := tgbotapi.NewMessage(chatId, "抽卡记录同步成功！")
	sendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(sendMessage)
	return true, nil
}
