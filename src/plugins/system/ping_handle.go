package system

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
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
	go utils.DelayDelMsg(msg.Chat.ID, msg.MessageID, viper.GetDuration("bot.msg_del_delay"))
	return true, nil
}
