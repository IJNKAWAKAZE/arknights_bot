package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/cache"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
)

func ClearHandle(update tgbotapi.Update) error {
	owner := viper.GetInt64("bot.owner")
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	param := update.Message.CommandArguments()
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if param == "" {
		sendMessage := tgbotapi.NewMessage(chatId, "参数不能为空")
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	if owner == userId {
		res, ctx := cache.RedisScanKeys(param)
		for res.Next(ctx) {
			cache.RedisDel(res.Val())
		}
		sendMessage := tgbotapi.NewMessage(chatId, "清理成功")
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := config.Arknights.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}
