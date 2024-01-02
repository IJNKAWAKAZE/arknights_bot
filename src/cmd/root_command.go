package cmd

import (
	"arknights_bot/config"
	"arknights_bot/core/bot"
	"arknights_bot/core/cron"
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
	//初始化redis连接
	config.Redis()
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
	//开始消息监听
	bot.Serve()
}
