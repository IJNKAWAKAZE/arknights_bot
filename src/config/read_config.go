package config

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
)

var MsgDelDelay float64
var HeadhuntTimes int
var PoolUP = make(map[int]string)
var Pool = make(map[int]string)

var RecruitMissing map[string]string
var RecruitTagMap = make(map[string]string)

func init() {
	// 设置配置文件的名字
	viper.SetConfigName("arknights")
	// 设置配置文件的类型
	viper.SetConfigType("yaml")
	// 添加配置文件的路径
	viper.AddConfigPath("./")
	// 寻找配置文件并读取
	err := viper.ReadInConfig()
	if err != nil {
		log.Println(err)
		return
	}
	initData()
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed")
		initData()
	})
}

func initData() {
	MsgDelDelay = viper.GetFloat64("bot.msg_del_delay")
	HeadhuntTimes = viper.GetInt("headhunt.times")
	PoolUP[7] = viper.GetString("headhunt.pool_up_6_1")
	PoolUP[6] = viper.GetString("headhunt.pool_up_6")
	PoolUP[5] = viper.GetString("headhunt.pool_up_5")
	Pool[6] = viper.GetString("headhunt.pool_6")
	Pool[5] = viper.GetString("headhunt.pool_5")
	Pool[4] = viper.GetString("headhunt.pool_4")
	Pool[3] = viper.GetString("headhunt.pool_3")
	jpMissing := viper.GetString("recruit.missing.jp")
	recruitTags := viper.GetString("recruit.tags")
	RecruitMissing = make(map[string]string)
	RecruitTagMap = make(map[string]string)
	for _, missing := range strings.Split(jpMissing, "/") {
		RecruitMissing[missing] = missing
	}
	for _, tag := range strings.Split(recruitTags, "/") {
		t := strings.Split(tag, "-")
		RecruitTagMap[t[0]] = t[1]
	}
}
