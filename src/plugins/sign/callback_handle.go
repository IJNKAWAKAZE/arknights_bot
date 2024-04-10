package sign

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

// SignPlayer 选择签到角色
func SignPlayer(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 3 {
		return nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	clickUserId, _ := strconv.ParseInt(d[1], 10, 64)
	uid := d[2]

	if userId != clickUserId {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的角色！")
		bot.Arknights.Send(answer)
		return nil
	}

	var userAccount account.UserAccount
	var player account.UserPlayer

	utils.GetAccountByUid(userId, uid).Scan(&userAccount)
	utils.GetPlayerByUserId(userId, uid).Scan(&player)

	answer := tgbotapi.NewCallback(callbackQuery.ID, "")
	bot.Arknights.Send(answer)
	return Sign(player, userAccount, chatId)
}
