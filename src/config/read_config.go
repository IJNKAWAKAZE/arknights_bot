package config

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"log"
	"strings"
	"sync"
)

var MsgDelDelay float64
var HeadhuntTimes int

// DataMu 保护以下配置数据的并发读写（viper 热更新时写入）
var DataMu sync.RWMutex

var PoolUP = make(map[int]string)
var Pool = make(map[int]string)
var RecruitMissing map[string]string
var RecruitTagMap map[string]string
var EnemyName map[string]string
var IgnoreBirthday map[string]string
var ADWords []string

func init() {
	// 包初始化阶段加载默认配置，失败仅记录日志；启动时由 LoadConfig 显式校验并报错退出
	// 设置配置文件的名字
	viper.SetConfigName("arknights")
	// 设置配置文件的类型
	viper.SetConfigType("yaml")
	// 添加配置文件的路径
	viper.AddConfigPath("./")
	// 寻找配置文件并读取
	if err := viper.ReadInConfig(); err != nil {
		log.Println(err)
		return
	}
	initData()
	watchConfig()
}

// LoadConfig 显式加载配置文件，path 为空时使用默认路径（./arknights.yaml）
func LoadConfig(path string) error {
	if path != "" {
		viper.SetConfigFile(path)
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
		initData()
		watchConfig()
		return nil
	}
	if viper.ConfigFileUsed() == "" {
		return errors.New("未找到配置文件 arknights.yaml，请复制 arknights.example.yaml 并修改")
	}
	return nil
}

func watchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed")
		initData()
	})
}

func initData() {
	DataMu.Lock()
	defer DataMu.Unlock()
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
	enemyName := viper.GetString("enemy_name")
	RecruitMissing = make(map[string]string)
	RecruitTagMap = make(map[string]string)
	EnemyName = make(map[string]string)
	IgnoreBirthday = make(map[string]string)
	for _, ignore := range viper.GetStringSlice("ignore_birthday") {
		IgnoreBirthday[ignore] = ignore
	}
	ADWords = viper.GetStringSlice("ad")
	for _, missing := range strings.Split(jpMissing, "/") {
		RecruitMissing[missing] = missing
	}
	if len(recruitTags) > 0 {
		for _, tag := range strings.Split(recruitTags, "/") {
			t := strings.Split(tag, "-")
			RecruitTagMap[t[0]] = t[1]
		}
	}
	if len(enemyName) > 0 {
		for _, enemy := range strings.Split(enemyName, "/") {
			t := strings.Split(enemy, "-")
			EnemyName[t[0]] = t[1]
		}
	}
}
