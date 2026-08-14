package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/media"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
)

// CalendarHandle 活动日历
func CalendarHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID

	_, _ = config.Arknights.SendChatAction(chatId, "upload_photo")
	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/calendar", port), 0, 1.5)
	if err != nil {
		msg, _ := config.Arknights.ReplyText(chatId, messageId, err.Error())
		messagecleaner.AddDelQueue(chatId, msg.MessageID, 5)
		return nil
	}
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "calendar.jpg"})
	sendDocument.ReplyToMessageID = messageId
	config.Arknights.Send(sendDocument)
	return nil
}
