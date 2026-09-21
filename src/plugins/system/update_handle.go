package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/datasource"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/localassets"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
)

func UpdateHandle(update tgbotapi.Update) error {
	owner := viper.GetInt64("bot.owner")
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if owner == userId {
		msg, err := config.Arknights.ReplyText(chatId, messageId, "开始更新数据源")
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		if err := datasource.UpdateDataSourceRunner(); err != nil {
			msg, sendErr := config.Arknights.SendText(chatId, "数据源更新失败："+err.Error()+"，详见日志")
			if sendErr != nil {
				return sendErr
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
			return nil
		}
		text := "数据源更新结束"
		if localassets.Enabled() {
			text += "\n本地素材正在后台增量同步，首次全量下载较慢且会占用较多磁盘空间，完成情况见日志。"
		}
		msg, err = config.Arknights.SendText(chatId, text)
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
