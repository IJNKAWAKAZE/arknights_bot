package config

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
)

var Arknights *tgbotapi.BotAPI

func Bot() error {
	token := viper.GetString("bot.token")
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Println(err)
		return err
	}
	Arknights = bot
	log.Println("机器人初始化完成")
	return nil
}
