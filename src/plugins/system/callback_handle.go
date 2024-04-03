package system

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	messageId := callbackQuery.Message.MessageID
	target, _ := strconv.ParseInt(d[2], 10, 64)
	targetMessageId, _ := strconv.Atoi(d[3])

	if !utils.IsAdmin(chatId, userId) {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "无使用权限！")
		bot.Arknights.Send(answer)
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
		delMsg = tgbotapi.NewDeleteMessage(chatId, messageId)
		bot.Arknights.Send(delMsg)
	}

	if d[1] == "CLOSE" {
		delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
		bot.Arknights.Send(delMsg)
	}
	answer := tgbotapi.NewCallback(callbackQuery.ID, "")
	bot.Arknights.Send(answer)
	return nil
}
