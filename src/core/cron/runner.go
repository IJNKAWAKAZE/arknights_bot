package cron

import (
	"arknights_bot/plugins/arknightsnews"
	"arknights_bot/plugins/datasource"
	"github.com/robfig/cron/v3"
	"log"
)

func StartCron() error {
	crontab := cron.New(cron.WithSeconds())

	//明日方舟bilibili动态 0 0/10 * * * ?
	_, err := crontab.AddFunc("0 0/10 * * * ?", arknightsnews.BilibiliNews())
	if err != nil {
		return err
	}

	//每周五凌晨2点33更新数据源 0 33 02 ? * FRI
	_, err = crontab.AddFunc("0 33 02 ? * FRI", datasource.UpdateDataSource())
	if err != nil {
		return err
	}

	//启动定时任务
	crontab.Start()
	log.Println("定时任务已启动")
	return nil
}
