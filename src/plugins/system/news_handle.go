package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func NewsHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if bot.Arknights.IsAdmin(chatId, userId) {
		var joined utils.GroupJoined
		utils.GetJoinedByChatId(chatId).Scan(&joined)
		joined.News = joined.News ^ 1
		bot.DBEngine.Table("group_joined").Save(&joined)
		text := "动态推送已开启！"
		if joined.News == 0 {
			text = "动态推送已关闭！"
		}
		sendMessage := tgbotapi.NewMessage(chatId, text)
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
	sendMessage.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendMessage)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	return nil
}
