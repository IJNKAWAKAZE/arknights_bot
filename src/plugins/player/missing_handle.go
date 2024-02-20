package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandOperation"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"regexp"
)

// MissingHandle 未获取干员
type PlayerOperationMissing struct {
	commandOperation.OperationAbstract
}

// BoxHandle 我的干员

func (_ PlayerOperationMissing) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) (bool, error) {

	sendAction := tgbotapi.NewChatAction(chatId, "upload_document")
	bot.Arknights.Send(sendAction)
	messageId := message.MessageID
	param := message.CommandArguments()
	matched, _ := regexp.MatchString("^[0-9\\d]+(,[0-9\\d]+)*$", param)
	if param != "" && param != "all" && !matched {
		sendMessage := tgbotapi.NewMessage(chatId, "参数错误")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/missing?userId=%d&uid=%s&param=%s", port, userAccount.UserNumber, uid, param), 0)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败，token可能已失效请重设token。")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}

	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "missing.png"})
	sendDocument.ReplyToMessageID = messageId
	bot.Arknights.Send(sendDocument)
	return true, nil
}
