package player

import (
	"arknights_bot/plugins/account"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// PlayerHandle 角色信息查询
func PlayerHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

	var userAccount account.UserAccount
	var players []account.UserPlayer

	userAccountI, playersI, err := getAccountAndPlayers(update)
	if err != nil && userAccountI != nil && playersI != nil {
		userAccount = *userAccountI
		players = *playersI
	} else {
		return true, nil
	}
	operation, ok := parseStringToOperation(update.Message.Command())
	if !ok {
		return true, nil
	}
	if len(players) > 1 {
		return true, playerSelector(chatId, userId, messageId, players, operation)

	}
	// 判断操作类型
	switch operation {
	case OP_STATE:
		return State(players[0].Uid, userAccount, chatId, messageId)
	case OP_BOX:
		param := update.Message.CommandArguments()
		return Box(players[0].Uid, userAccount, chatId, messageId, param)
	case OP_GACHA:
		return Gacha(players[0].Uid, userAccount, chatId, messageId)
	case OP_CARD:
		return Card(players[0].Uid, userAccount, chatId, messageId)
	case OP_EXPORT:
		return Export(players[0].Uid, userAccount, chatId)
	}
	return true, nil
}
