package config

import (
	"github.com/spf13/viper"
	"log"
)

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
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}

func GetInt(key string) int {
	return viper.GetInt(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}
