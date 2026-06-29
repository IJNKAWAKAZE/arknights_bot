package cmd

import (
	"arknights_bot/config"
	"arknights_bot/core/bot"
	"arknights_bot/core/cron"
	"arknights_bot/core/web"
	"arknights_bot/plugins/antispam"
)

func Execute() {
	Launch()
}

func Launch() {
	//初始化数据库连接
	err := config.DB()
	if err != nil {
		panic(err)
	}
	// 数据库迁移要求手动执行，启动时不要自动改表结构。
	// err = config.MigrateDB()
	// if err != nil {
	// 	panic(err)
	// }
	//初始化redis连接
	config.Redis()
	antispam.Init()
	//初始化机器人
	err = config.Bot()
	if err != nil {
		panic(err)
	}
	//开启定时任务
	err = cron.StartCron()
	if err != nil {
		panic(err)
	}
	//开启http服务
	go web.Start()
	//开始消息监听
	bot.Serve()
}
