package system

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// HelpHandle 帮助
func HelpHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	port := viper.GetString("http.port")
	pic := utils.Screenshot("http://localhost:" + port + "/help")
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	bot.Arknights.Send(sendPhoto)
	return true, nil
}
