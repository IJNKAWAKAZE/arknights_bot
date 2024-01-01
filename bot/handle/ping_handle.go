package handle

import (
	bot "arknights_bot/bot/init"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func PingHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	sendSticker := tgbotapi.NewSticker(chatId, tgbotapi.FileID("CAACAgUAAxkBAAPBZYTpA2TK3cnNG1CfXiLG_9Yo7t8AAiYEAAML4FddSAYCSfeu2jME"))
	sendSticker.ReplyToMessageID = messageId
	bot.Arknights.Send(sendSticker)
	return true, nil
}
