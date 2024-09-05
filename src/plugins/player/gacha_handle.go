package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/spf13/viper"
	"log"
	"time"
)

// GachaHandle 抽卡记录
type PlayerOperationGacha struct {
	commandoperation.OperationAbstract
}

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

func (_ PlayerOperationGacha) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	messageId := message.MessageID
	token := userAccount.HypergryphToken
	channelId := "1"
	var userPlayer account.UserPlayer
	utils.GetPlayerByUserId(userAccount.UserNumber, uid).Scan(&userPlayer)
	if userPlayer.ServerName == "b服" || userPlayer.ServerName == "bilibili服" {
		token = userPlayer.BToken
		channelId = "2"
		// BToken为空设置BToken
		if token == "" {
			sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("BToken未设置，请先进行[设置](https://t.me/%s)。", viper.GetString("bot.name")))
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			sendMessage.ReplyToMessageID = messageId
			msg, err := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(chatId, messageId, 5)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}
	}
	// 获取角色抽卡记录
	chars, err := skland.GetPlayerGacha(token, channelId)
	if err != nil {
		log.Println(err)
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return err
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
				UserName:   userAccount.UserName,
				UserNumber: userAccount.UserNumber,
				Uid:        uid,
				PoolName:   c.PoolName,
				PoolOrder:  c.PoolOrder,
				CharName:   c.Name,
				IsNew:      c.IsNew,
				Rarity:     c.Rarity,
				Ts:         c.Ts,
			}
			bot.DBEngine.Table("user_gacha").Create(&userGacha)
		}
	}

	var userGacha []UserGacha
	res := utils.GetUserGacha(userAccount.UserNumber, uid).Scan(&userGacha)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "不存在抽卡记录。")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	pic, e := utils.Screenshot(fmt.Sprintf("http://localhost:%s/gacha?userId=%d&uid=%s", port, userAccount.UserNumber, uid), 3000, 1.5)
	if e != nil {
		sendMessage := tgbotapi.NewMessage(chatId, e.Error())
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}

	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "gacha.jpg"})
	sendDocument.ReplyToMessageID = messageId
	bot.Arknights.Send(sendDocument)
	return nil
}
