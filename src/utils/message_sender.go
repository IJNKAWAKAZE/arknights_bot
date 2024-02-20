package utils

import (
	"arknights_bot/config"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func SendMessage(chatId int64, massage string, isMarkDown bool, replyId *int) {
	tgMassage := tgbotapi.NewMessage(chatId, massage)
	if replyId != nil {
		tgMassage.ReplyToMessageID = *replyId
	}
	if isMarkDown {
		tgMassage.ParseMode = tgbotapi.ModeMarkdownV2
	}
	msg, err := config.Arknights.Send(tgMassage)
	if err != nil {
		log.Printf("%v can not be send error : %v", msg, err)
	}

}
