package main

import (
	"arknights_bot/config"
	"arknights_bot/plugins/cron"
	"arknights_bot/plugins/gatekeeper"
	"log"
)

func main() {
	//启动
	Launch()
}

// 设置日志格式
func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func Launch() {
	//初始化数据库连接
	config.DB()
	//初始化redis连接
	config.Redis()
	//初始化机器人
	config.Bot()
	//开启定时任务
	cron.StartCron()
	//开始消息监听
	gatekeeper.Processor()
}
