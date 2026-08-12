package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"strconv"
	"strings"
)

func RegulationHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	messagecleaner.AddDelQueue(chatId, messageId, 5)

	if config.Arknights.IsAdmin(chatId, userId) {
		replyToMessage := update.Message.ReplyToMessage
		if replyToMessage != nil {
			replyMessageId := replyToMessage.MessageID
			var joined model.GroupJoined
			repo.GetJoinedByChatId(chatId).Scan(&joined)
			joined.Reg = replyMessageId
			config.DBEngine.Table("group_joined").Save(&joined)
			var sendMessage tgbotapi.MessageConfig
			if replyToMessage.Chat.UserName != "" {
				sendMessage = tgbotapi.NewMessage(chatId, fmt.Sprintf("消息[%d](https://t.me/%s/%d)已设置为群规！", replyMessageId, replyToMessage.Chat.UserName, replyMessageId))
			} else {
				sendMessage = tgbotapi.NewMessage(chatId, fmt.Sprintf("消息[%d](https://t.me/c/%s/%d)已设置为群规！", replyMessageId, strings.ReplaceAll(strconv.FormatInt(replyToMessage.Chat.ID, 10), "-100", ""), replyMessageId))
			}
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
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
