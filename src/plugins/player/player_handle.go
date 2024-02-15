package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// PlayerHandle 角色信息查询
func PlayerHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID

	var userAccount account.UserAccount
	var players []account.UserPlayer

	userAccountP, playersP, err := getAccountAndPlayers(update)
	if err != nil || userAccountP == nil || playersP == nil {
		return true, err
	}
	operations, ok := parseStringToOperation(update.Message.Command())
	if !ok {
		return true, nil
	}
	players = *playersP
	userAccount = *userAccountP
	if len(players) > 1 {
		if operations != OP_REDEEM {
			return true, playerSelector(update, players, operations, NO_REQUIREMENT)
		}
		cdk := strings.ToUpper(update.Message.CommandArguments())
		if cdk == "" {
			SendMessage := tgbotapi.NewMessage(chatId, "请输入CDK！")
			SendMessage.ReplyToMessageID = messageId
			bot.Arknights.Send(SendMessage)
			return true, nil
		}

		Redeem[userId] = cdk
	}
	// 判断操作类型
	switch operations {
	case OP_EXPORT:
		return Export(players[0].Uid, userAccount, chatId)
	case OP_STATE:
		return State(players[0].Uid, userAccount, chatId, messageId)
	case OP_BOX:
		param := update.Message.CommandArguments()
		return Box(players[0].Uid, userAccount, chatId, messageId, param)
	case OP_MISSING:
		param := update.Message.CommandArguments()
		return Missing(players[0].Uid, userAccount, chatId, messageId, param)
	case OP_SYNC:
		return SyncGacha(players[0].Uid, userAccount, chatId, messageId, userAccount.UserName)
	case OP_GACHA:
		return Gacha(players[0].Uid, userAccount, chatId, messageId)
	case OP_CARD:
		return Card(players[0].Uid, userAccount, chatId, messageId)
	case OP_BASE:
		return Base(players[0].Uid, userAccount, chatId, messageId)
	case OP_REDEEM:
		cdk := strings.ToUpper(update.Message.CommandArguments())
		return RedeemCDK(players[0].Uid, userAccount, chatId, messageId, cdk)
	}
	return true, nil
}
