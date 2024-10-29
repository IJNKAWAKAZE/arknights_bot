package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"strings"
)

func NewMemberHandle(update tgbotapi.Update) error {
	message := update.Message
	for _, member := range message.NewChatMembers {
		if member.ID == message.From.ID { // 自己加入群组
			chat, err := bot.Arknights.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: member.ID}})
			if err != nil {
				return err
			}
			for _, word := range bot.ADWords {
				if strings.Contains(chat.Bio, word) {
					message.Delete()
					bot.Arknights.BanChatMember(message.Chat.ID, member.ID)
					return nil
				}
			}
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
