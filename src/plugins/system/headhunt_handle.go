package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

var PoolUP = make(map[int]string)
var Pool = make(map[int]string)

func init() {
	PoolUP[6] = viper.GetString("gacha.pool_up_6")
	PoolUP[5] = viper.GetString("gacha.pool_up_5")
	Pool[6] = viper.GetString("gacha.pool_6")
	Pool[5] = viper.GetString("gacha.pool_5")
	Pool[4] = viper.GetString("gacha.pool_4")
	Pool[3] = viper.GetString("gacha.pool_3")
}

// HeadhuntHandle 寻访模拟
func HeadhuntHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 60)
	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)
	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/headhunt?userId=%d", port, userId), 0)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "失败啦，啊哈哈哈！")
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 5)
		return true, nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	msg, _ := bot.Arknights.Send(sendPhoto)
	messagecleaner.AddDelQueue(chatId, msg.MessageID, 60)
	return true, nil
}
