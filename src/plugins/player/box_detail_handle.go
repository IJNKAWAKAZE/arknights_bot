package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"regexp"
)

type PlayerOperationBoxDetail struct {
	commandoperation.OperationAbstract
}

// BoxHandle 我的干员

func (_ PlayerOperationBoxDetail) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	messageId := message.MessageID
	param := message.CommandArguments()
	sendAction := tgbotapi.NewChatAction(chatId, "upload_document")
	bot.Arknights.Send(sendAction)

	matched, _ := regexp.MatchString("^[1-6]$", param)
	if param != "" && !matched {
		sendMessage := tgbotapi.NewMessage(chatId, "参数错误")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}

	port := viper.GetString("http.port")
	pic, err := utils.Screenshot(fmt.Sprintf("http://localhost:%s/boxDetail?userId=%d&uid=%s&param=%s&sklandId=%s", port, userAccount.UserNumber, uid, param, userAccount.SklandId), 0, 1.5)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}

	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "box_detail.jpg"})
	sendDocument.ReplyToMessageID = messageId
	bot.Arknights.Send(sendDocument)
	return nil
}
