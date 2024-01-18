package bot

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/gatekeeper"
	"arknights_bot/plugins/operator"
	"arknights_bot/plugins/player"
	"arknights_bot/plugins/sign"
	"arknights_bot/plugins/system"
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
	bot.TeleBot.NewCallBackProcessor("verify", gatekeeper.CallBackData)
	bot.TeleBot.NewCallBackProcessor("bind", account.ChoosePlayer)
	bot.TeleBot.NewCallBackProcessor("unbind", account.UnbindPlayer)
	bot.TeleBot.NewCallBackProcessor("sign", sign.SignPlayer)
	bot.TeleBot.NewCallBackProcessor("player", player.PlayerData)
	bot.TeleBot.NewCallBackProcessor("report", system.Report)

	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.Message != nil && len(update.Message.NewChatMembers) > 0
	}, gatekeeper.NewMemberHandle)
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.Message != nil && update.Message.LeftChatMember != nil
	}, gatekeeper.LeftMemberHandle)
	bot.TeleBot.NewProcessor(func(update tgbotapi.Update) bool {
		return update.InlineQuery != nil
	}, gatekeeper.InlineQueryHandle)
	// 私聊
	bot.TeleBot.NewPrivateCommandProcessor("start", system.HelpHandle)
	bot.TeleBot.NewPrivateCommandProcessor("cancel", account.CancelHandle)
	bot.TeleBot.NewPrivateCommandProcessor("bind", account.BindHandle)
	bot.TeleBot.NewPrivateCommandProcessor("unbind", account.UnbindHandle)
	bot.TeleBot.NewPrivateCommandProcessor("reset_token", account.SetTokenHandle)
	bot.TeleBot.NewPrivateCommandProcessor("sync_gacha", player.SyncGachaHandle)

	bot.TeleBot.NewWaitMessageProcessor("setToken", account.SetToken)
	bot.TeleBot.NewWaitMessageProcessor("resetToken", account.ResetToken)

	// 普通
	bot.TeleBot.NewCommandProcessor("help", system.HelpHandle)
	bot.TeleBot.NewCommandProcessor("ping", system.PingHandle)
	bot.TeleBot.NewCommandProcessor("sign", sign.SignHandle)
	bot.TeleBot.NewCommandProcessor("state", player.PlayerHandle)
	bot.TeleBot.NewCommandProcessor("box", player.PlayerHandle)
	bot.TeleBot.NewCommandProcessor("card", player.PlayerHandle)
	bot.TeleBot.NewCommandProcessor("gacha", player.PlayerHandle)
	bot.TeleBot.NewCommandProcessor("operator", operator.OperatorHandle)
	bot.TeleBot.NewCommandProcessor("report", system.ReportHandle)

	// 权限
	bot.TeleBot.NewCommandProcessor("update", system.UpdateHandle)
	bot.TeleBot.Run(bot.Arknights.GetUpdatesChan(u))
}
