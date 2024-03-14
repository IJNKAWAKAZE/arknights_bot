package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HelpCmd struct {
	PrivateCmds []Cmd
	PublicCmds  []Cmd
	AdminCmds   []Cmd
}

type Cmd struct {
	Cmd    string
	Desc   string
	Param  string
	IsBind bool
}

func Help(r *gin.Engine) {
	r.GET("/help", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Help.tmpl")
		var helpCmd HelpCmd
		var privateCmds []Cmd
		var publicCmds []Cmd
		var adminCmds []Cmd

		// 私聊指令
		privateCmds = append(privateCmds, Cmd{Cmd: "/bind", Desc: "绑定角色", Param: "", IsBind: false})
		privateCmds = append(privateCmds, Cmd{Cmd: "/unbind", Desc: "解绑角色", Param: "", IsBind: true})
		privateCmds = append(privateCmds, Cmd{Cmd: "/resume", Desc: "设置名片签名", Param: "", IsBind: true})
		privateCmds = append(privateCmds, Cmd{Cmd: "/cancel", Desc: "取消操作", Param: "", IsBind: false})
		privateCmds = append(privateCmds, Cmd{Cmd: "/reset_token", Desc: "重设token", Param: "", IsBind: true})
		privateCmds = append(privateCmds, Cmd{Cmd: "/btoken", Desc: "设置BToken", Param: "", IsBind: true})
		privateCmds = append(privateCmds, Cmd{Cmd: "/sync_gacha", Desc: "同步抽卡记录", Param: "", IsBind: true})
		privateCmds = append(privateCmds, Cmd{Cmd: "/import_gacha", Desc: "导入抽卡记录", Param: "", IsBind: true})
		privateCmds = append(privateCmds, Cmd{Cmd: "/export_gacha", Desc: "导出抽卡记录", Param: "", IsBind: true})

		// 普通指令
		publicCmds = append(publicCmds, Cmd{Cmd: "/help", Desc: "使用说明", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/ping", Desc: "存活测试", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/sign", Desc: "签到", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/sign", Desc: "开启自动签到", Param: "auto", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/sign", Desc: "关闭自动签到", Param: "stop", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/state", Desc: "当前状态", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/box", Desc: "我的干员(默认6星)", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/box", Desc: "所有干员", Param: "all", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/box", Desc: "对应星级干员", Param: "5,6", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/missing", Desc: "未获取干员(默认6星)", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/missing", Desc: "所有未获取干员", Param: "all", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/missing", Desc: "对应星级未获取干员", Param: "5,6", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/card", Desc: "我的名片", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/base", Desc: "基建信息", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/gacha", Desc: "抽卡记录", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/operator", Desc: "干员查询", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/enemy", Desc: "敌人查询", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/report", Desc: "举报", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/quiz", Desc: "云玩家检测", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/quiz", Desc: "云玩家检测(困难)", Param: "h", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/redeem", Desc: "CDK兑换", Param: "[CDK]", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/headhunt", Desc: "寻访模拟", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/recruit", Desc: "公招计算(图片附带)", Param: "", IsBind: false})

		// 管理员指令
		adminCmds = append(adminCmds, Cmd{Cmd: "/update", Desc: "更新数据源", Param: "", IsBind: false})
		adminCmds = append(adminCmds, Cmd{Cmd: "/news", Desc: "开启/关闭动态推送", Param: "", IsBind: false})
		adminCmds = append(adminCmds, Cmd{Cmd: "/quiz", Desc: "开启云玩家检测", Param: "start", IsBind: false})
		adminCmds = append(adminCmds, Cmd{Cmd: "/quiz", Desc: "关闭云玩家检测", Param: "stop", IsBind: false})
		adminCmds = append(adminCmds, Cmd{Cmd: "/headhunt", Desc: "开启寻访模拟", Param: "start", IsBind: false})
		adminCmds = append(adminCmds, Cmd{Cmd: "/headhunt", Desc: "关闭寻访模拟", Param: "stop", IsBind: false})

		helpCmd.PrivateCmds = privateCmds
		helpCmd.PublicCmds = publicCmds
		helpCmd.AdminCmds = adminCmds
		c.HTML(http.StatusOK, "Help.tmpl", helpCmd)
	})
}
