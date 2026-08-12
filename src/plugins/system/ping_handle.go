package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
)

// PingHandle 存活测试
func PingHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	sendSticker := tgbotapi.NewSticker(chatId, tgbotapi.FileID(viper.GetString("sticker.ping")))
	sendSticker.ReplyToMessageID = messageId
	msg, err := config.Arknights.Send(sendSticker)
	messagecleaner.AddDelQueue(chatId, messageId, 5)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}
