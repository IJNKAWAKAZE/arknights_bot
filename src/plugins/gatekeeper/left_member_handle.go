package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func LeftMemberHandle(update tgbotapi.Update) error {
	update.Message.Delete()
	// 清除离开成员的验证状态，避免残留条目影响后续加入
	if update.Message.LeftChatMember != nil {
		verifySet.checkExistAndRemove(update.Message.LeftChatMember.ID, update.Message.Chat.ID)
	}
	return nil
}
