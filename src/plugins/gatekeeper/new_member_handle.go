package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func NewMemberHandle(update tgbotapi.Update) (bool, error) {
	message := update.Message
	delMsg := tgbotapi.NewDeleteMessage(message.Chat.ID, message.MessageID)
	bot.Arknights.Send(delMsg)
	for _, member := range message.NewChatMembers {
		if member.ID == message.From.ID { // 自己加入群组
			go VerifyMember(message)
			continue
		}
		// 邀请加入群组，无需进行验证
		utils.SaveInvite(message, &member)
		name := utils.GetFullName(message.From)
		newName := utils.GetFullName(&member)
		sendMessage := tgbotapi.NewMessage(message.Chat.ID,
			fmt.Sprintf("[%s](tg://user?id=%d)邀请了[%s](tg://user?id=%d)加入群组。",
				name, message.From.ID, newName, member.ID))
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		msg, _ := bot.Arknights.Send(sendMessage)
		go utils.DelayDelMsg(msg.Chat.ID, msg.MessageID, viper.GetDuration("bot.msg_del_delay"))
	}
	return true, nil
}
