package cmd

import (
	"arknights_bot/config"
	"arknights_bot/core/bot"
	"arknights_bot/core/cron"
	"arknights_bot/core/shutdown"
	"arknights_bot/core/web"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func Execute() {
	configPath := flag.String("config", "", "配置文件路径（默认 ./arknights.yaml）")
	flag.Parse()
	if err := config.LoadConfig(*configPath); err != nil {
		log.Fatalf("加载配置文件失败：%v", err)
	}
	Launch()
}

func Launch() {
	//初始化数据库连接
	err := config.DB()
	if err != nil {
		panic(err)
	}
	//初始化redis连接
	if err = config.Redis(); err != nil {
		log.Fatalf("Redis连接失败：%v", err)
	}
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
	//注册优雅退出清理
	shutdown.Register(cron.Stop)
	shutdown.Register(web.Shutdown)
	shutdown.Register(config.Close)
	//开启http服务
	go web.Start()
	//监听退出信号
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("收到退出信号")
		shutdown.All()
	}()
	//开始消息监听
	bot.Serve()
}
