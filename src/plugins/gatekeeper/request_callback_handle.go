package gatekeeper

import (
	bot "arknights_bot/config"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"strconv"
	"strings"
)

func RequestCallBackData(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return nil
	}

	userId, _ := strconv.ParseInt(d[1], 10, 64)
	chatId, _ := strconv.ParseInt(d[2], 10, 64)

	if has, correct := verifySet.checkExistAndRemove(userId, chatId); has {
		if d[3] != correct {
			callbackQuery.Answer(true, "验证未通过")
			declineChatJoinRequest := tgbotapi.DeclineChatJoinRequest{ChatConfig: tgbotapi.ChatConfig{ChatID: chatId}, UserID: userId}
			bot.Arknights.Request(declineChatJoinRequest)
		} else {
			callbackQuery.Answer(true, "验证通过！")
			approveChatJoinRequest := tgbotapi.ApproveChatJoinRequestConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: chatId}, UserID: userId}
			bot.Arknights.Request(approveChatJoinRequest)
		}
		callbackQuery.Delete()
	}
	return nil
}
