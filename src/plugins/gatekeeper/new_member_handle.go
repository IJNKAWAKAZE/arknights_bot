package gatekeeper

import (
	"arknights_bot/config"
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
)

func NewMemberHandle(update tgbotapi.Update) error {
	message := update.Message
	var joined model.GroupJoined
	repo.GetJoinedByChatId(message.Chat.ID).Scan(&joined)
	if joined.RequestMode == 1 { // 不使用此验证
		// 入群审计：审批模式下进群，应来自通过入群验证后的批准
		for _, member := range message.NewChatMembers {
			if member.UserName == viper.GetString("bot.name") {
				continue
			}
			auditJoin(message.Chat.ID, member.ID, member.FullName(), "审批模式进群", "入群申请被批准后进群")
		}
		return nil
	}
	for _, member := range message.NewChatMembers {
		chatId := message.Chat.ID
		userId := member.ID
		if member.ID == message.From.ID { // 自己加入群组
			verifySet.add(userId, chatId, "")
			auditJoin(chatId, userId, member.FullName(), "自行进群", "进入人工验证流程")
			// 昵称广告词检查
			if hasAdWord(member.FullName()) {
				log.Printf("入群验证：用户 %d 昵称含广告词，踢出", userId)
				auditJoin(chatId, userId, member.FullName(), "踢出", "昵称含广告词")
				message.Delete()
				config.Arknights.BanChatMember(chatId, userId)
				verifySet.checkExistAndRemove(userId, chatId)
				return nil
			}
			chat, err := config.Arknights.GetChatInfo(member.ID)
			if err != nil {
				log.Println("获取用户信息失败", err)
			} else if hasAdWord(chat.Bio) {
				log.Printf("入群验证：用户 %d bio 含广告词，踢出", userId)
				auditJoin(chatId, userId, member.FullName(), "踢出", "bio 含广告词")
				message.Delete()
				config.Arknights.BanChatMember(chatId, userId)
				verifySet.checkExistAndRemove(userId, chatId)
				return nil
			}
			go VerifyMember(message)
			continue
		}
		// 机器人被邀请加群
		if member.UserName == viper.GetString("bot.name") {
			repo.SaveJoined(message)
			continue
		}
		// 邀请加入群组，无需进行验证
		auditJoin(chatId, userId, member.FullName(), "邀请入群", "被邀请加入，无需验证")
		repo.SaveInvite(message, &member)
	}
	return nil
}
