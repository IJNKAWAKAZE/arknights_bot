package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
)

type PlayerOperationBoxSummary struct {
	commandoperation.OperationAbstract
}

// BoxHandle 我的干员

func (_ PlayerOperationBoxSummary) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	messageId := message.MessageID
	sendAction := tgbotapi.NewChatAction(chatId, "upload_document")
	bot.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	pic, err := utils.Screenshot(fmt.Sprintf("http://localhost:%s/boxSummary?userId=%d&uid=%s&sklandId=%s", port, userAccount.UserNumber, uid, userAccount.SklandId), 0, 1.5)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}

	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "box_summary.jpg"})
	sendDocument.ReplyToMessageID = messageId
	bot.Arknights.Send(sendDocument)
	return nil
}
