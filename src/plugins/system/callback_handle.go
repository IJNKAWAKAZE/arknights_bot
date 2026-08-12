package system

import (
	"arknights_bot/config"
	"fmt"
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

	if !config.Arknights.IsAdminWithPermissions(chatId, userId, tgbotapi.AdminCanRestrictMembers) {
		callbackQuery.Answer(true, "无使用权限！")
		return nil
	}

	if d[1] == "BAN" {
		fmt.Printf("在群组 %s 用户 %s 封禁了 %d", callbackQuery.Message.Chat.Title, callbackQuery.From.FullName(), target)
		config.Arknights.BanChatMember(chatId, target)
		delMsg := tgbotapi.NewDeleteMessage(chatId, targetMessageId)
		config.Arknights.Send(delMsg)

	}
	callbackQuery.Delete()
	callbackQuery.Answer(false, "")
	return nil
}
