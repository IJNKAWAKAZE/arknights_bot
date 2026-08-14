package bot

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/apremind"
	"arknights_bot/plugins/enemy"
	"arknights_bot/plugins/gatekeeper"
	"arknights_bot/plugins/lottery"
	"arknights_bot/plugins/operator"
	"arknights_bot/plugins/player"
	"arknights_bot/plugins/sign"
	"arknights_bot/plugins/skin"
	"arknights_bot/plugins/system"
	"github.com/spf13/viper"
	"log"
)

// Serve TG机器人阻塞监听
func Serve() {
	log.Println("机器人启动成功")
	b := config.Arknights
	b.Debug = viper.GetBool("bot.debug")
	b.IgnoreChannelCMD = true
	b.JoinRequestProcessor(gatekeeper.JoinRequestHandle)
	b.NewMemberProcessor(gatekeeper.NewMemberHandle)
	b.LeftMemberProcessor(gatekeeper.LeftMemberHandle)

	// callback
	b.NewCallBackProcessor("verify", gatekeeper.CallBackData)
	b.NewCallBackProcessor("request", gatekeeper.RequestCallBackData)
	b.NewCallBackProcessor("chooseServer", account.ChooseServer)
	b.NewCallBackProcessor("bind", account.ChoosePlayer)
	b.NewCallBackProcessor("unbind", account.UnbindPlayer)
	b.NewCallBackProcessor("sign", sign.SignPlayer)
	b.NewCallBackProcessor("player", player.PlayerData)
	b.NewCallBackProcessor("report", system.Report)

	// InlineQuery
	b.NewInlineQueryProcessor("干员", operator.InlineOperator)
	b.NewInlineQueryProcessor("皮肤", skin.InlineSkin)
	b.NewInlineQueryProcessor("敌人", enemy.InlineEnemy)

	// 私聊
	b.NewPrivateCommandProcessor("start", system.HelpHandle)
	b.NewPrivateCommandProcessor("cancel", account.CancelHandle)
	b.NewPrivateCommandProcessor("bind", account.BindHandle)
	b.NewPrivateCommandProcessor("unbind", account.UnbindHandle)
	b.NewPrivateCommandProcessor("reset_token", account.SetTokenHandle)
	b.NewPrivateCommandProcessor("import_gacha", player.PlayerHandle)
	b.NewPrivateCommandProcessor("export_gacha", player.PlayerHandle)

	// wait
	b.NewWaitMessageProcessor("setToken", account.SetToken)
	b.NewWaitMessageProcessor("resetToken", account.ResetToken)
	b.NewWaitMessageProcessor("importGacha", player.PlayerHandle)

	// 普通
	b.NewCommandProcessor("help", system.HelpHandle)
	b.NewCommandProcessor("ping", system.PingHandle)
	b.NewCommandProcessor("tag", system.TagHandle)
	b.NewCommandProcessor("sign", sign.SignHandle)
	b.NewCommandProcessor("ap", apremind.ApHandle)
	b.NewCommandProcessor("state", player.PlayerHandle)
	b.NewCommandProcessor("box", player.PlayerHandle)
	b.NewCommandProcessor("box_detail", player.PlayerHandle)
	b.NewCommandProcessor("box_summary", player.PlayerHandle)
	b.NewCommandProcessor("missing", player.PlayerHandle)
	b.NewCommandProcessor("card", player.PlayerHandle)
	b.NewCommandProcessor("base", player.PlayerHandle)
	b.NewCommandProcessor("gacha", player.PlayerHandle)
	b.NewCommandProcessor("operator", operator.OperatorHandle)
	b.NewCommandProcessor("skin", skin.SkinHandle)
	b.NewCommandProcessor("enemy", enemy.EnemyHandle)
	b.NewCommandProcessor("report", system.ReportHandle)
	b.NewCommandProcessor("quiz", system.QuizHandle)
	b.NewCommandProcessor("redeem", player.PlayerHandle)
	b.NewCommandProcessor("headhunt", system.HeadhuntHandle)
	b.NewCommandProcessor("calendar", system.CalendarHandle)
	b.NewCommandProcessor("depot", player.PlayerHandle)
	b.NewCommandProcessor("join_lottery", lottery.JoinLotteryHandle)
	b.NewCommandProcessor("lottery_detail", lottery.LotteryDetailHandle)

	// 图片
	b.NewPhotoMessageProcessor("/recruit", system.RecruitHandle)
	//回复
	b.NewReplyMessageProcessor("/recruit", system.ReplyRecruitHandle)

	// 管理员
	b.NewAdminCommandProcessor("news", system.NewsHandle)
	b.NewAdminCommandProcessor("birthday", system.BirthdayHandle)
	b.NewAdminCommandProcessor("request_mode", system.RequestModeHandle)
	b.NewAdminCommandProcessor("reg", system.RegulationHandle)
	b.NewAdminCommandProcessor("welcome", system.WelcomeHandle)
	b.NewAdminCommandProcessor("start_lottery", lottery.StartLotteryHandle)
	b.NewAdminCommandProcessor("stop_lottery", lottery.StopLotteryHandle)
	b.NewAdminCommandProcessor("end_lottery", lottery.EndLotteryHandle)
	b.NewAdminCommandProcessor("lottery", lottery.LotteryHandle)

	// 仅拥有者
	b.NewOwnerCommandProcessor("update", system.UpdateHandle)
	b.NewOwnerCommandProcessor("sign_all", sign.SignAllHandle)
	b.NewOwnerCommandProcessor("clear", system.ClearHandle)
	b.NewOwnerCommandProcessor("kill", system.KillHandle)
	b.Run()
}
