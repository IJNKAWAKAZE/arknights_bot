package cron

import (
	"github.com/robfig/cron/v3"
	"log"
)

func StartCron() {
	crontab := cron.New(cron.WithSeconds())

	//明日方舟bilibili动态 0 0/10 * * * ?
	crontab.AddFunc("0 0/10 * * * ?", BilibiliNews())

	//每周五凌晨2点33更新数据源 0 33 02 ? * FRI
	crontab.AddFunc("0 33 02 ? * FRI", UpdateDataSource())

	//清理消息 0 0/1 * * * ?
	crontab.AddFunc("0 0/1 * * * ?", DelMsg())

	//启动定时任务
	crontab.Start()
	log.Println("定时任务已启动")
}
