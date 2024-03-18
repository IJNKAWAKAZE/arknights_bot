package gatekeeper

import (
	bot "arknights_bot/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func LeftMemberHandle(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	messageId := message.MessageID
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)
	return nil
}
