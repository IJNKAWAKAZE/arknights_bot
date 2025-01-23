package gatekeeper

import (
	bot "arknights_bot/config"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func CheckMember(update tgbotapi.Update) bool {
	if update.Message != nil && update.Message.Text != "" && verifySet.checkExist(update.SentFrom().ID, update.FromChat().ID) {
		return true
	}
	return false
}

func KickMember(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	message.Delete()
	bot.Arknights.BanChatMember(chatId, userId)
	verifyC <- nil
	return nil
}
