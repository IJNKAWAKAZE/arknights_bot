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
	var joined utils.GroupJoined
	utils.GetJoinedByChatId(message.Chat.ID).Scan(&joined)
	for _, member := range message.NewChatMembers {
		chatId := message.Chat.ID
		userId := member.ID
		if member.ID == message.From.ID { // 自己加入群组
			// 清除可能残留的验证状态（如申请模式下管理员手动放行遗留的条目）
			verifySet.checkExistAndRemove(userId, chatId)
			if joined.RequestMode == 1 && !recentlyApproved.is(userId, chatId) {
				// 请求模式下未经答题进入（如未开启 Telegram 审核时通过链接加入），仍需验证
				verifySet.add(userId, chatId, "")
				chat, err := bot.Arknights.GetChatInfo(member.ID)
				if err != nil {
					return err
				}
				for _, word := range bot.ADWords {
					if strings.Contains(chat.Bio, word) {
						message.Delete()
						bot.Arknights.BanChatMember(chatId, userId)
						return nil
					}
				}
				go VerifyMember(message)
				continue
			}
			if joined.RequestMode == 1 {
				// 请求模式下已由机器人批准（答对验证），无需重复验证
				continue
			}
			verifySet.add(userId, chatId, "")
			chat, err := bot.Arknights.GetChatInfo(member.ID)
			if err != nil {
				return err
			}
			for _, word := range bot.ADWords {
				if strings.Contains(chat.Bio, word) {
					message.Delete()
					bot.Arknights.BanChatMember(chatId, userId)
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
		// 邀请加入群组，无需进行验证（信任管理员邀请）；清除其可能残留的验证状态
		verifySet.checkExistAndRemove(userId, chatId)
		utils.SaveInvite(message, &member)
	}
	return nil
}
