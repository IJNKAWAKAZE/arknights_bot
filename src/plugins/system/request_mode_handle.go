package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func RequestModeHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if config.Arknights.IsAdmin(chatId, userId) {
		var joined model.GroupJoined
		if err := repo.CheckDB(repo.GetJoinedByChatId(chatId).Scan(&joined), chatId); err != nil {
			return err
		}
		joined.RequestMode = joined.RequestMode ^ 1
		if err := repo.CheckDB(config.DBEngine.Table("group_joined").Save(&joined), chatId); err != nil {
			return err
		}
		text := "请求模式开启！"
		if joined.RequestMode == 0 {
			text = "请求模式关闭！"
		}
		msg, err := config.Arknights.ReplyText(chatId, messageId, text)
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
