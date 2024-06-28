package system

import (
	bot "arknights_bot/config"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"strconv"
	"strings"
)

// Report 举报
func Report(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	target, _ := strconv.ParseInt(d[2], 10, 64)
	targetMessageId, _ := strconv.Atoi(d[3])

	if !bot.Arknights.IsAdminWithPermissions(chatId, userId, 16) {
		callbackQuery.Answer(true, "无使用权限！")
		return nil
	}

	if d[1] == "BAN" {
		banChatMemberConfig := tgbotapi.BanChatMemberConfig{
			ChatMemberConfig: tgbotapi.ChatMemberConfig{
				ChatID: chatId,
				UserID: target,
			},
			RevokeMessages: true,
		}
		bot.Arknights.Send(banChatMemberConfig)
		delMsg := tgbotapi.NewDeleteMessage(chatId, targetMessageId)
		bot.Arknights.Send(delMsg)

	}
	callbackQuery.Delete()
	callbackQuery.Answer(false, "")
	return nil
}
