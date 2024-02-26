package system

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
)

var FileId string

// HelpHandle 帮助
func HelpHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	if FileId == "" {
		port := viper.GetString("http.port")
		pic := utils.Screenshot("http://localhost:"+port+"/help", 0, 1.5)
		if pic == nil {
			sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败！")
			sendMessage.ReplyToMessageID = messageId
			bot.Arknights.Send(sendMessage)
			return true, nil
		}
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
		sendPhoto.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendPhoto)
		if err != nil {
			log.Println(err)
			return true, err
		}
		FileId = msg.Photo[0].FileID
		return true, nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(FileId))
	sendPhoto.ReplyToMessageID = messageId
	bot.Arknights.Send(sendPhoto)
	return true, nil
}
