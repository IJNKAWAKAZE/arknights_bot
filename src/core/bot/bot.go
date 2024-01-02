package bot

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/gatekeeper"
	"arknights_bot/utils/telebot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

// Serve TG机器人阻塞监听
func Serve() {
	log.Println("机器人启动成功")
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	// 注册处理器
	bot.TeleBot = &telebot.Bot{}
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.CallbackQuery != nil
	}, gatekeeper.CallBackData)
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.Message != nil && len(update.Message.NewChatMembers) > 0
	}, gatekeeper.NewMemberHandle)
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.Message != nil && update.Message.LeftChatMember != nil
	}, gatekeeper.LeftMemberHandle)
	bot.TeleBot.NewCommandProcessor("ping", gatekeeper.PingHandle)

	bot.TeleBot.Run(bot.Arknights.GetUpdatesChan(u))
}
