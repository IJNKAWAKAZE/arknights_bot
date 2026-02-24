package gatekeeper

import (
	"arknights_bot/utils"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func JoinRequestHandle(update tgbotapi.Update) error {
	var joined utils.GroupJoined
	utils.GetJoinedByChatId(update.ChatJoinRequest.Chat.ID).Scan(&joined)
	if joined.RequestMode == 0 { // 不使用此验证
		return nil
	}
	go VerifyRequestMember(update)
	return nil
}
