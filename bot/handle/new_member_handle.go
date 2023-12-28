package handle

import (
	"arknights_bot/bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func NewMemberHandle(update tgbotapi.Update) (bool, error) {
	message := update.Message
	chatId := message.Chat.ID
	userId := message.From.ID
	messageId := message.MessageID
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	utils.DeleteMessage(delMsg)
	for _, member := range message.NewChatMembers {
		memberId := member.ID
		if memberId == userId {
			go Verify(message)
		} else {
			utils.SaveInvite(message, member)
			name := utils.GetFullName(message.From)
			newName := utils.GetNewMemberName(member)
			sendMessage := tgbotapi.NewMessage(chatId, "<a href=\"tg://user?id="+strconv.FormatInt(userId, 10)+"\">"+name+"</a>邀请了<a href=\"tg://user?id="+strconv.FormatInt(memberId, 10)+"\">"+newName+"</a>加入群组。")
			sendMessage.ParseMode = tgbotapi.ModeHTML
			msg, _ := utils.SendMessage(sendMessage)
			utils.AddDelQueue(msg.Chat.ID, msg.MessageID, 2)
		}
	}
	return true, nil
}
