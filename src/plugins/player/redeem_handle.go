package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"strings"
)

type PlayerOperationRedeem struct {
	commandoperation.OperationAbstract
}

// RedeemHandle CDK兑换
func (_ PlayerOperationRedeem) CheckRequirementsAndPrepare(update tgbotapi.Update) bool {
	return len(update.Message.CommandArguments()) != 0
}
func (_ PlayerOperationRedeem) HintOnRequirementsFailed() (string, bool) {
	return "请输入CDK！", false
}
func (_ PlayerOperationRedeem) HintWordForPlayerSelection() string {
	return "请选择要兑换的角色"
}
func (_ PlayerOperationRedeem) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	messageId := message.MessageID
	cdk := message.CommandArguments()
	cdk = strings.ToUpper(cdk)
	token := userAccount.HypergryphToken
	channelId := "1"
	var userPlayer account.UserPlayer
	utils.GetPlayerByUserId(userAccount.UserNumber, uid).Scan(&userPlayer)
	if userPlayer.ServerName == "b服" || userPlayer.ServerName == "bilibili服" {
		token = userPlayer.BToken
		channelId = "2"
		// BToken为空设置BToken
		if token == "" {
			sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("BToken未设置，请先进行[设置](https://t.me/%s)。", viper.GetString("bot.name")))
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			sendMessage.ReplyToMessageID = messageId
			msg, err := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(chatId, messageId, 5)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}
	}
	result, err := skland.GetPlayerRedeem(token, cdk, channelId)
	if err != nil {
		SendMessage := tgbotapi.NewMessage(chatId, err.Error())
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
		return err
	}
	if result != "" {
		SendMessage := tgbotapi.NewMessage(chatId, result)
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
		return fmt.Errorf(result)
	}
	SendMessage := tgbotapi.NewMessage(chatId, "CDK兑换成功，请进入游戏领取奖励。")
	SendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(SendMessage)
	return nil
}
