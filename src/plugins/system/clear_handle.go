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
		msg, err := config.Arknights.ReplyText(chatId, messageId, "参数不能为空")
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
		msg, err := config.Arknights.ReplyText(chatId, messageId, "清理成功")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	msg, err := config.Arknights.ReplyText(chatId, messageId, "无使用权限！")
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}
