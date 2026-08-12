package gatekeeper

import (
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

func JoinRequestHandle(update tgbotapi.Update) error {
	var joined model.GroupJoined
	repo.GetJoinedByChatId(update.ChatJoinRequest.Chat.ID).Scan(&joined)
	if joined.RequestMode == 0 { // 不使用此验证
		return nil
	}
	go VerifyRequestMember(update)
	return nil
}
