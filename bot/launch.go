package bot

import (
	"arknights_bot/bot/handle"
	initConfig "arknights_bot/bot/init"
	"arknights_bot/bot/scheduled"
	"log"
)

// 设置日志格式
func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
}

func Launch() {
	//初始化数据库连接
	initConfig.DB()
	//初始化redis连接
	initConfig.Redis()
	//初始化机器人
	initConfig.Bot()
	//开启定时任务
	scheduled.StartCron()
	//开始消息监听
	handle.Processor()
}
