package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func WelcomeHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if config.Arknights.IsAdmin(chatId, userId) {
		text := update.Message.CommandArguments()
		if text != "" {
			var joined model.GroupJoined
			if err := repo.CheckDB(repo.GetJoinedByChatId(chatId).Scan(&joined), chatId); err != nil {
				return err
			}
			joined.Welcome = text
			if err := repo.CheckDB(config.DBEngine.Table("group_joined").Save(&joined), chatId); err != nil {
				return err
			}
			msg, err := config.Arknights.SendText(chatId, "设置入群欢迎信息成功")
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		}
		return nil
	}

	msg, err := config.Arknights.ReplyText(chatId, messageId, "无使用权限！")
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}
