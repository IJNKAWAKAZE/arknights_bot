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

// RedeemHandle CDK兑换
func RedeemHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
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

	res := utils.GetAccountByUserId(userId).Scan(&userAccount)
	if res.RowsAffected == 0 {
		// 未绑定账号
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("未查询到绑定账号，请先进行[绑定](https://t.me/%s)。", viper.GetString("bot.name")))
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return true, nil
	}

	// 获取绑定角色
	res = utils.GetPlayersByUserId(userId).Scan(&players)
	if res.RowsAffected == 0 {
		sendMessage := tgbotapi.NewMessage(chatId, "您还未绑定任何角色！")
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, 5)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return true, nil
	}

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
