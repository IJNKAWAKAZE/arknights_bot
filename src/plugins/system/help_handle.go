package system

import (
	"arknights_bot/config"
	"arknights_bot/utils/media"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
)

var fileId string

// HelpHandle 帮助
func HelpHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	config.Arknights.Send(sendAction)

	if fileId == "" {
		port := viper.GetString("http.port")
		pic, err := media.Screenshot("http://localhost:"+port+"/help", 0, 1.5)
		if err != nil {
			sendMessage := tgbotapi.NewMessage(chatId, err.Error())
			sendMessage.ReplyToMessageID = messageId
			config.Arknights.Send(sendMessage)
			return nil
		}
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
		sendPhoto.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendPhoto)
		if err != nil {
			log.Println(err)
			return err
		}
		fileId = msg.Photo[0].FileID
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(fileId))
	sendPhoto.ReplyToMessageID = messageId
	config.Arknights.Send(sendPhoto)
	return nil
}
