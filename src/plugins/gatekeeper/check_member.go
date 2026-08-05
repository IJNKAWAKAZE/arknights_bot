package gatekeeper

import (
	bot "arknights_bot/config"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

// CheckMember 匹配验证期间抢先发送消息的用户
func CheckMember(update tgbotapi.Update) bool {
	if update.Message != nil && update.Message.Text != "" && verifySet.checkExist(update.SentFrom().ID, update.FromChat().ID) {
		return true
	}
	return false
}

// KickMember 踢出验证期间抢发消息的用户；仅对仍处于验证中的用户生效，
// 避免误伤已通过验证或被管理员手动放行的用户
func KickMember(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	message.Delete()
	if has, _ := verifySet.checkExistAndRemove(userId, chatId); !has {
		return nil
	}
	bot.Arknights.BanChatMember(chatId, userId)
	return nil
}
