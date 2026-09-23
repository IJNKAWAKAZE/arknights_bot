package config

import (
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

var Arknights *tgbotapi.Bot

func Bot() error {
	token := viper.GetString("bot.token")
	bot, err := tgbotapi.NewBotAPIWithClient(token, tgbotapi.APIEndpoint, &http.Client{Timeout: 75 * time.Second})
	if err != nil {
		log.Println(err)
		return err
	}
	Arknights = bot.AddHandle()
	Arknights.SetOwnerID(viper.GetInt64("bot.owner"))
	log.Println("机器人初始化完成")
	return nil
}
