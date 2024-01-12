package gatekeeper

import (
	"arknights_bot/plugins/operator"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func InlineQueryHandle(update tgbotapi.Update) (bool, error) {
	inlineQuery := update.InlineQuery
	operator.InlineOperator(inlineQuery)
	return true, nil
}
