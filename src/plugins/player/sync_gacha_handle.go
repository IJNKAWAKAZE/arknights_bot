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
func SyncGachaHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

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

	// 遍历角色
	for _, player := range players {
		if player.ServerName == "b服" {
			SendMessage := tgbotapi.NewMessage(chatId, "暂不支持B服！")
			bot.Arknights.Send(SendMessage)
		}

		// 获取角色抽卡记录
		chars, err := skland.GetPlayerGacha(userAccount.HypergryphToken)
		if err != nil {
			log.Println(err)
			sendMessage := tgbotapi.NewMessage(chatId, "token可能已失效请重设token。")
			sendMessage.ReplyToMessageID = messageId
			bot.Arknights.Send(sendMessage)
			return true, err
		}

		// 获取上次更新时间
		var lastUpdate int64
		bot.DBEngine.Raw("select ts from user_gacha where user_number = ? and uid = ? order by ts desc limit 1", player.UserNumber, player.Uid).Scan(&lastUpdate)

		// 同步抽卡数据
		for _, c := range chars {
			if c.Ts > lastUpdate {
				id, _ := gonanoid.New(32)
				userGacha := UserGacha{
					Id:         id,
					UserName:   utils.GetFullName(update.Message.From),
					UserNumber: userId,
					Uid:        player.Uid,
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
	}

	sendMessage := tgbotapi.NewMessage(chatId, "抽卡记录同步成功！")
	sendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(sendMessage)
	return true, nil
}
