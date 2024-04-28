package gatekeeper

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func LeftMemberHandle(update tgbotapi.Update) error {
	update.Message.Delete()
	return nil
}
