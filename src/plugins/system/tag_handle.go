package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func TagHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	param := update.Message.CommandArguments()
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, config.MsgDelDelay)
	if config.Arknights.IsAdmin(chatId, userId) {
		sendMessage := tgbotapi.NewMessage(chatId, "此指令仅普通成员可使用！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	if param == "" {
		sendMessage := tgbotapi.NewMessage(chatId, "请输入要设置的标签！")
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	_, err := config.Arknights.SetMemberTag(chatId, userId, param)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return err
	}
	sendMessage := tgbotapi.NewMessage(chatId, "自定义标签成功！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := config.Arknights.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
	return nil
}
