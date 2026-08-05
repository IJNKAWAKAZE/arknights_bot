package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func JoinRequestHandle(update tgbotapi.Update) error {
	var joined utils.GroupJoined
	utils.GetJoinedByChatId(update.ChatJoinRequest.Chat.ID).Scan(&joined)
	if joined.RequestMode == 0 { // 不使用此验证
		return nil
	}
	chatId := update.ChatJoinRequest.Chat.ID
	userId := update.ChatJoinRequest.From.ID
	if limiter.isBlocked(chatId, userId) {
		// 多次答题失败被封禁的用户，直接拒绝申请
		bot.Arknights.DeclineChatJoinRequest(chatId, userId)
		return nil
	}
	go VerifyRequestMember(update)
	return nil
}
