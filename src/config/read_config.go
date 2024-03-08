package config

import (
	"github.com/spf13/viper"
	"log"
)

var MsgDelDelay float64
var HeadhuntTimes int

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
	MsgDelDelay = viper.GetFloat64("bot.msg_del_delay")
	HeadhuntTimes = viper.GetInt("headhunt.times")
}
