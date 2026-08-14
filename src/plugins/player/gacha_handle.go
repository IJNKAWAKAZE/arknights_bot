package player

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/media"
	"arknights_bot/utils/repo"
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
	if userAccount.ServerName == "国际服" {
		config.Arknights.ReplyText(chatId, messageId, "国际服暂不可用")
		return nil
	}
	token := userAccount.HypergryphToken
	// 获取角色抽卡记录
	chars, err := skland.GetPlayerGacha(token, uid)
	if err != nil {
		log.Println(err)
		config.Arknights.SendMarkdownV2(chatId, err.Error(), messageId)
		return err
	}

	// 获取上次更新时间
	var lastUpdate int64
	config.DBEngine.Raw("select ts from user_gacha where user_number = ? and uid = ? order by ts desc limit 1", userAccount.UserNumber, uid).Scan(&lastUpdate)

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
			config.DBEngine.Table("user_gacha").Create(&userGacha)
		}
	}

	var userGacha []UserGacha
	res := repo.GetUserGacha(userAccount.UserNumber, uid).Scan(&userGacha)
	if res.RowsAffected == 0 {
		config.Arknights.ReplyText(chatId, messageId, "不存在抽卡记录。")
		return nil
	}

	_, _ = config.Arknights.SendChatAction(chatId, "upload_photo")

	port := viper.GetString("http.port")
	pic, e := media.Screenshot(fmt.Sprintf("http://localhost:%s/gacha?userId=%d&uid=%s", port, userAccount.UserNumber, uid), 3000, 1.5)
	if e != nil {
		config.Arknights.ReplyText(chatId, messageId, e.Error())
		return nil
	}

	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "gacha.jpg"})
	sendDocument.ReplyToMessageID = messageId
	config.Arknights.Send(sendDocument)
	return nil
}
