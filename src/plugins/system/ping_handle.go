package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

// PingHandle 存活测试
func PingHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	sendSticker := tgbotapi.NewSticker(chatId, tgbotapi.FileID(viper.GetString("sticker.ping")))
	sendSticker.ReplyToMessageID = messageId
	msg, _ := bot.Arknights.Send(sendSticker)
	messagecleaner.AddDelQueue(chatId, messageId, 5)
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return true, nil
}
