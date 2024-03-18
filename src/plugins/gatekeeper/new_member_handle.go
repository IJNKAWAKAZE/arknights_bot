package gatekeeper

import (
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func NewMemberHandle(update tgbotapi.Update) error {
	message := update.Message
	for _, member := range message.NewChatMembers {
		if member.ID == message.From.ID { // 自己加入群组
			go VerifyMember(message)
			continue
		}
		// 机器人被邀请加群
		if member.UserName == viper.GetString("bot.name") {
			utils.SaveJoined(message)
			continue
		}
		// 邀请加入群组，无需进行验证
		utils.SaveInvite(message, &member)
	}
	return nil
}
