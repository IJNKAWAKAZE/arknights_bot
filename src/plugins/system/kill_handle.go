package system

import (
	"arknights_bot/config"
	"arknights_bot/core/shutdown"
	"arknights_bot/plugins/messagecleaner"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
)

func KillHandle(update tgbotapi.Update) error {
	owner := viper.GetInt64("bot.owner")
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if owner == userId {
		log.Println("关闭机器人")
		shutdown.All()
		return nil
	}

	msg, err := config.Arknights.ReplyText(chatId, messageId, "无使用权限！")
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}
