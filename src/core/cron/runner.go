package cron

import (
	"arknights_bot/plugins/arknightsnews"
	"arknights_bot/plugins/datasource"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/sign"
	"arknights_bot/plugins/system"
	"github.com/robfig/cron/v3"
	"log"
)

func StartCron() error {
	crontab := cron.New(cron.WithSeconds())

	//明日方舟bilibili动态 0/30 * * * * ?
	_, err := crontab.AddFunc("0/30 * * * * ?", arknightsnews.BilibiliNews)
	if err != nil {
		return err
	}

	//每周五凌晨2点33更新数据源 0 33 02 ? * FRI
	_, err = crontab.AddFunc("0 33 02 ? * FRI", datasource.UpdateDataSource())
	if err != nil {
		return err
	}

	//每日1点执行自动签到 0 0 1 * * ?
	_, err = crontab.AddFunc("0 0 1 * * ?", sign.AutoSign())
	if err != nil {
		return err
	}

	//清理消息 0/1 * * * * ?
	_, err = crontab.AddFunc("0/1 * * * * ?", messagecleaner.DelMsg)
	if err != nil {
		return err
	}

	//重置每日寻访次数 0 0 0 * * ?
	_, err = crontab.AddFunc("0 0 0 * * ?", system.ResetHeadhuntTimes())
	if err != nil {
		return err
	}

	//启动定时任务
	crontab.Start()
	log.Println("定时任务已启动")
	return nil
}
