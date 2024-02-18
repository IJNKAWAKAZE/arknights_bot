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
	"strings"
)

var redeem = make(map[int64]string)

// RedeemHandle CDK兑换
func getRedeemPerFeq(update tgbotapi.Update) bool {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	messageId := update.Message.MessageID
	cdk := strings.ToUpper(update.Message.CommandArguments())
	if cdk == "" {
		SendMessage := tgbotapi.NewMessage(chatId, "请输入CDK！")
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
		return false
	} else {
		redeem[userId] = cdk
		return true
	}
}

func RedeemCDK(uid string, userAccount account.UserAccount, chatId int64, messageId int, cdk string) (bool, error) {
	delete(redeem, userAccount.UserNumber)
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
