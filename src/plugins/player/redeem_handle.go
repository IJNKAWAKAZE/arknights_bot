package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

// RedeemHandle CDK兑换
func RedeemHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	cdk := strings.ToUpper(update.Message.CommandArguments())

	if cdk == "" {
		SendMessage := tgbotapi.NewMessage(chatId, "请输入CDK！")
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
		return true, nil
	}

	var userAccount account.UserAccount
	var players []account.UserPlayer

	userAccountI, playersI, err := getAccountAndPlayers(update)
	if err != nil {
		return true, err
	} else if playersI == nil || userAccountI == nil {
		return true, nil
	}
	userAccount = *userAccountI
	players = *playersI

	// 遍历角色
	for _, player := range players {
		if player.ServerName == "b服" {
			SendMessage := tgbotapi.NewMessage(chatId, "暂不支持B服！")
			SendMessage.ReplyToMessageID = messageId
			bot.Arknights.Send(SendMessage)
			return true, nil
		}
		result, err := skland.GetPlayerRedeem(userAccount.HypergryphToken, cdk)
		if err != nil {
			SendMessage := tgbotapi.NewMessage(chatId, err.Error())
			SendMessage.ReplyToMessageID = messageId
			bot.Arknights.Send(SendMessage)
			return true, err
		}
		if result != "" {
			SendMessage := tgbotapi.NewMessage(chatId, result)
			SendMessage.ReplyToMessageID = messageId
			bot.Arknights.Send(SendMessage)
			return true, fmt.Errorf(result)
		}
		SendMessage := tgbotapi.NewMessage(chatId, "CDK兑换成功，请进入游戏领取奖励。")
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
	}
	return true, nil
}
