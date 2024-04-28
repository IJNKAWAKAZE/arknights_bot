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

// MissingHandle 未获取干员
type PlayerOperationMissing struct {
	commandoperation.OperationAbstract
}

func (_ PlayerOperationMissing) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {

	sendAction := tgbotapi.NewChatAction(chatId, "upload_document")
	bot.Arknights.Send(sendAction)
	messageId := message.MessageID
	param := message.CommandArguments()
	matched, _ := regexp.MatchString("^[0-9\\d]+(,[0-9\\d]+)*$", param)
	if param != "" && param != "all" && !matched {
		sendMessage := tgbotapi.NewMessage(chatId, "参数错误")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/missing?userId=%d&uid=%s&param=%s&sklandId=%s", port, userAccount.UserNumber, uid, param, userAccount.SklandId), 0, 1.5)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("生成图片失败，token可能已失效请[重设token](https://t.me/%s)。", viper.GetString("bot.name")))
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}

	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "missing.jpg"})
	sendDocument.ReplyToMessageID = messageId
	bot.Arknights.Send(sendDocument)
	return nil
}
