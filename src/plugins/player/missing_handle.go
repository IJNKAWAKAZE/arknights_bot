package player

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils/media"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"regexp"
)

// MissingHandle 未获取干员
type PlayerOperationMissing struct {
	commandoperation.OperationAbstract
}

func (_ PlayerOperationMissing) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	messageId := message.MessageID
	if userAccount.ServerName == "国际服" {
		config.Arknights.ReplyText(chatId, messageId, "国际服暂不可用")
		return nil
	}
	_, _ = config.Arknights.SendChatAction(chatId, "upload_document")
	param := message.CommandArguments()
	matched, _ := regexp.MatchString("^[1-6](,[1-6])*$", param)
	if param != "" && param != "all" && !matched {
		config.Arknights.ReplyText(chatId, messageId, "参数错误")
		return nil
	}

	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/missing?userId=%d&uid=%s&param=%s&sklandId=%s", port, userAccount.UserNumber, uid, param, userAccount.SklandId), 0, 1.5)
	if err != nil {
		config.Arknights.SendMarkdownV2(chatId, err.Error(), messageId)
		return nil
	}

	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "missing.jpg"})
	sendDocument.ReplyToMessageID = messageId
	config.Arknights.Send(sendDocument)
	return nil
}
