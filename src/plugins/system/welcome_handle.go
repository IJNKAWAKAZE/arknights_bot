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
			repo.GetJoinedByChatId(chatId).Scan(&joined)
			joined.Welcome = text
			config.DBEngine.Table("group_joined").Save(&joined)
			sendMessage := tgbotapi.NewMessage(chatId, "设置入群欢迎信息成功")
			msg, err := config.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		}
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
