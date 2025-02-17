package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"strings"
	"time"
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
	if utils.RedisIsExists("risk_control") {
		SendMessage := tgbotapi.NewMessage(chatId, "触发风控，请等待解除后再进行兑换！")
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
		return nil
	}
	token := userAccount.HypergryphToken
	result, err := skland.GetPlayerRedeem(token, cdk, uid)
	if err != nil {
		SendMessage := tgbotapi.NewMessage(chatId, err.Error())
		SendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(SendMessage)
		return err
	}
	if result != "" {
		if result == "需要验证" {
			utils.RedisSet("risk_control", "1", time.Hour)
		}
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
