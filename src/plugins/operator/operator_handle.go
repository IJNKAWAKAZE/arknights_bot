package operator

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// OperatorHandle 干员查询
func OperatorHandle(update tgbotapi.Update) (bool, error) {
	text := ""
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text:                         "选择干员",
				SwitchInlineQueryCurrentChat: &text,
			},
		),
	)
	sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "请选择要查询的干员")
	sendMessage.ReplyMarkup = inlineKeyboardMarkup
	msg, _ := bot.Arknights.Send(sendMessage)
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return true, nil
}
