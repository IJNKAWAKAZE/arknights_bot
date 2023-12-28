package handle

import (
	bot "arknights_bot/bot/init"
	"arknights_bot/pkg/telebot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func init() {
	bot.TeleBot = &telebot.Bot{}
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.CallbackQuery != nil
	}, CallBackData)
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.Message != nil && len(update.Message.NewChatMembers) > 0
	}, NewMemberHandle)
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.Message != nil && update.Message.LeftChatMember != nil
	}, LeftMemberHandle)
}

func Processor() {
	log.Println("机器人启动成功")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := bot.Kawakaze.GetUpdatesChan(u)
	bot.TeleBot.Run(updates)
}
