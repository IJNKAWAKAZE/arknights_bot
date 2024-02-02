package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

var Redeem = make(map[int64]string)

// RedeemHandle CDK兑换
func RedeemHandle(players []account.UserPlayer, userAccount account.UserAccount, chatId int64, userId int64, messageId int, cdk string) (bool, error) {
	if cdk == "" {
		SendMessage := tgbotapi.NewMessage(chatId, "请输入CDK！")
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
		return true, nil
	}

	Redeem[userId] = cdk
	if len(players) > 1 {
		// 绑定多个角色进行选择
		var buttons [][]tgbotapi.InlineKeyboardButton
		for _, player := range players {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(fmt.Sprintf("%s(%s)", player.PlayerName, player.ServerName), fmt.Sprintf("%s,%s,%d,%s,%d", "player", OP_REDEEM, userId, player.Uid, messageId)),
			))
		}
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要兑换的角色")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	} else {
		// 绑定单个角色
		return RedeemCDK(players[0].Uid, userAccount, chatId, messageId, cdk)
	}
	return true, nil
}

func RedeemCDK(uid string, userAccount account.UserAccount, chatId int64, messageId int, cdk string) (bool, error) {
	delete(Redeem, userAccount.UserNumber)
	token := userAccount.HypergryphToken
	channelId := "1"
	var userPlayer account.UserPlayer
	utils.GetPlayerByUserId(userAccount.UserNumber, uid).Scan(&userPlayer)
	if userPlayer.ServerName == "b服" || userPlayer.ServerName == "bilibili服" {
		token = userAccount.BToken
		channelId = "2"
		// BToken为空设置BToken
		if token == "" {
			sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("B服Token未设置，请先进行[设置](https://t.me/%s)。", viper.GetString("bot.name")))
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			sendMessage.ReplyToMessageID = messageId
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(chatId, messageId, 5)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return true, nil
		}
	}
	result, err := skland.GetPlayerRedeem(token, cdk, channelId)
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
	return true, nil
}
