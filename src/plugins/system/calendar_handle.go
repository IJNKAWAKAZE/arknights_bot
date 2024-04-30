package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
)

// CalendarHandle 活动日历
func CalendarHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)
	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/calendar", port), 0, 1)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败，请重试。")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 5)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	_, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		return err
	}
	return nil
}
