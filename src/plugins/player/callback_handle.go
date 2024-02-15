package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

// PlayerData 角色数据
func PlayerData(callBack tgbotapi.Update) (bool, error) {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 5 {
		return true, nil
	}

	userId := callbackQuery.From.ID
	chatId := callbackQuery.Message.Chat.ID
	operate, ok := parseIntStringToOperation(d[1])
	clickUserId, _ := strconv.ParseInt(d[2], 10, 64)
	uid := d[3]
	messageId, _ := strconv.Atoi(d[4])

	if userId != clickUserId {
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的角色！")
		bot.Arknights.Send(answer)
		return true, nil
	}
	if !ok {
		return true, nil
	}

	// 获取账号信息
	var userAccount account.UserAccount
	utils.GetAccountByUserId(userId).Scan(&userAccount)

	// 判断操作类型
	switch operate {
	case OP_STATE:
		return State(uid, userAccount, chatId, messageId)
	case OP_BOX:
		param := d[5]
		return Box(uid, userAccount, chatId, messageId, param)
	case OP_MISSING:
		param := d[5]
		return Missing(uid, userAccount, chatId, messageId, param)
	case OP_SYNC:
		name := utils.GetFullName(callbackQuery.From)
		return SyncGacha(uid, userAccount, chatId, messageId, name)
	case OP_GACHA:
		return Gacha(uid, userAccount, chatId, messageId)
	case OP_CARD:
		return Card(uid, userAccount, chatId, messageId)
	case OP_BASE:
		return Base(uid, userAccount, chatId, messageId)
	case OP_IMPORT:
		name := utils.GetFullName(callbackQuery.From)
		return Import(uid, userAccount, chatId, ImportFile[clickUserId], name)
	case OP_EXPORT:
		return Export(uid, userAccount, chatId)
	case OP_REDEEM:
		return RedeemCDK(uid, userAccount, chatId, messageId, Redeem[clickUserId])
	}

	return true, nil
}
