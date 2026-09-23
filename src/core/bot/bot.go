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
	b.JoinRequestProcessor(asyncControl(gatekeeper.JoinRequestHandle))
	b.NewMemberProcessor(asyncControl(gatekeeper.NewMemberHandle))
	b.LeftMemberProcessor(asyncControl(gatekeeper.LeftMemberHandle))

	// callback
	b.NewCallBackProcessor("verify", asyncControl(gatekeeper.CallBackData))
	b.NewCallBackProcessor("request", asyncControl(gatekeeper.RequestCallBackData))
	b.NewCallBackProcessor("chooseServer", async(account.ChooseServer))
	b.NewCallBackProcessor("bind", async(account.ChoosePlayer))
	b.NewCallBackProcessor("unbind", async(account.UnbindPlayer))
	b.NewCallBackProcessor("sign", async(sign.SignPlayer))
	b.NewCallBackProcessor("player", async(player.PlayerData))
	b.NewCallBackProcessor("report", async(system.Report))

	// InlineQuery
	b.NewInlineQueryProcessor("干员", async(operator.InlineOperator))
	b.NewInlineQueryProcessor("皮肤", async(skin.InlineSkin))
	b.NewInlineQueryProcessor("敌人", async(enemy.InlineEnemy))

	// 私聊
	b.NewPrivateCommandProcessor("start", async(system.HelpHandle))
	b.NewPrivateCommandProcessor("cancel", async(account.CancelHandle))
	b.NewPrivateCommandProcessor("bind", async(account.BindHandle))
	b.NewPrivateCommandProcessor("unbind", async(account.UnbindHandle))
	b.NewPrivateCommandProcessor("reset_token", async(account.SetTokenHandle))
	b.NewPrivateCommandProcessor("import_gacha", async(player.PlayerHandle))
	b.NewPrivateCommandProcessor("export_gacha", async(player.PlayerHandle))

	// wait
	b.NewWaitMessageProcessor("setToken", async(account.SetToken))
	b.NewWaitMessageProcessor("resetToken", async(account.ResetToken))
	b.NewWaitMessageProcessor("importGacha", async(player.PlayerHandle))

	// 普通
	b.NewCommandProcessor("help", async(system.HelpHandle))
	b.NewCommandProcessor("ping", async(system.PingHandle))
	b.NewCommandProcessor("tag", async(system.TagHandle))
	b.NewCommandProcessor("sign", async(sign.SignHandle))
	b.NewCommandProcessor("ap", async(apremind.ApHandle))
	b.NewCommandProcessor("state", async(player.PlayerHandle))
	b.NewCommandProcessor("box", async(player.PlayerHandle))
	b.NewCommandProcessor("box_detail", async(player.PlayerHandle))
	b.NewCommandProcessor("box_summary", async(player.PlayerHandle))
	b.NewCommandProcessor("missing", async(player.PlayerHandle))
	b.NewCommandProcessor("card", async(player.PlayerHandle))
	b.NewCommandProcessor("base", async(player.PlayerHandle))
	b.NewCommandProcessor("gacha", async(player.PlayerHandle))
	b.NewCommandProcessor("operator", async(operator.OperatorHandle))
	b.NewCommandProcessor("skin", async(skin.SkinHandle))
	b.NewCommandProcessor("enemy", async(enemy.EnemyHandle))
	b.NewCommandProcessor("report", async(system.ReportHandle))
	b.NewCommandProcessor("quiz", async(system.QuizHandle))
	b.NewCommandProcessor("redeem", async(player.PlayerHandle))
	b.NewCommandProcessor("headhunt", async(system.HeadhuntHandle))
	b.NewCommandProcessor("calendar", async(system.CalendarHandle))
	b.NewCommandProcessor("depot", async(player.PlayerHandle))
	b.NewCommandProcessor("join_lottery", async(lottery.JoinLotteryHandle))
	b.NewCommandProcessor("lottery_detail", async(lottery.LotteryDetailHandle))

	// 图片
	b.NewPhotoMessageProcessor("/recruit", async(system.RecruitHandle))
	//回复
	b.NewReplyMessageProcessor("/recruit", async(system.ReplyRecruitHandle))

	// 管理员
	b.NewAdminCommandProcessor("news", asyncControl(system.NewsHandle))
	b.NewAdminCommandProcessor("birthday", asyncControl(system.BirthdayHandle))
	b.NewAdminCommandProcessor("request_mode", asyncControl(system.RequestModeHandle))
	b.NewAdminCommandProcessor("reg", asyncControl(system.RegulationHandle))
	b.NewAdminCommandProcessor("welcome", asyncControl(system.WelcomeHandle))
	b.NewAdminCommandProcessor("start_lottery", async(lottery.StartLotteryHandle))
	b.NewAdminCommandProcessor("stop_lottery", async(lottery.StopLotteryHandle))
	b.NewAdminCommandProcessor("end_lottery", async(lottery.EndLotteryHandle))
	b.NewAdminCommandProcessor("lottery", async(lottery.LotteryHandle))

	// 仅拥有者
	b.NewOwnerCommandProcessor("update", async(system.UpdateHandle))
	b.NewOwnerCommandProcessor("sign_all", async(sign.SignAllHandle))
	b.NewOwnerCommandProcessor("clear", async(system.ClearHandle))
	b.NewOwnerCommandProcessor("kill", asyncControl(system.KillHandle))
	b.Run()
}
