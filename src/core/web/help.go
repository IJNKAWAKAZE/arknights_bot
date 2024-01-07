package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HelpCmd struct {
	PrivateCmds []Cmd
	PublicCmds  []Cmd
}

type Cmd struct {
	Cmd    string
	Desc   string
	Param  string
	IsBind bool
}

func Help(r *gin.Engine) {
	r.GET("/help", func(c *gin.Context) {
		var helpCmd HelpCmd
		var privateCmds []Cmd
		var publicCmds []Cmd

		// 私聊指令
		privateCmds = append(privateCmds, Cmd{Cmd: "/bind", Desc: "绑定角色", Param: "", IsBind: false})
		privateCmds = append(privateCmds, Cmd{Cmd: "/unbind", Desc: "解绑角色", Param: "", IsBind: true})
		privateCmds = append(privateCmds, Cmd{Cmd: "/cancel", Desc: "取消操作", Param: "", IsBind: false})

		// 普通指令
		publicCmds = append(publicCmds, Cmd{Cmd: "/help", Desc: "使用说明", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/ping", Desc: "存活测试", Param: "", IsBind: false})
		publicCmds = append(publicCmds, Cmd{Cmd: "/sign", Desc: "签到", Param: "", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/sign", Desc: "开启自动签到", Param: "auto", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/sign", Desc: "关闭自动签到", Param: "stop", IsBind: true})
		publicCmds = append(publicCmds, Cmd{Cmd: "/player", Desc: "当前状态", Param: "", IsBind: true})

		helpCmd.PrivateCmds = privateCmds
		helpCmd.PublicCmds = publicCmds
		c.HTML(http.StatusOK, "Help.tmpl", helpCmd)
	})
}
