package sign

import (
	"arknights_bot/plugins/account"
	"arknights_bot/utils"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
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
		callbackQuery.Answer(true, "这不是你的角色！")
		return nil
	}

	var userAccount account.UserAccount
	var player account.UserPlayer

	utils.GetAccountByUid(userId, uid).Scan(&userAccount)
	utils.GetPlayerByUserId(userId, uid).Scan(&player)

	callbackQuery.Answer(false, "")
	return Sign(player, userAccount, chatId)
}
