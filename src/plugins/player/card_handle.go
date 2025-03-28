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

// PlayerOperationCard 我的名片
type PlayerOperationCard struct {
	commandoperation.OperationAbstract
}

func (_ PlayerOperationCard) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	messageId := message.MessageID
	param := message.CommandArguments()
	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	r := "card"
	if param == "o" {
		r = "oldCard"
	}
	port := viper.GetString("http.port")
	pic, err := utils.Screenshot(fmt.Sprintf("http://localhost:%s/%s?userId=%d&uid=%s&sklandId=%s", port, r, userAccount.UserNumber, uid, userAccount.SklandId), 0, 1)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	/*sendPhoto.Caption = "点击复制UID:`" + uid + "`"
	sendPhoto.ParseMode = tgbotapi.ModeMarkdownV2*/
	bot.Arknights.Send(sendPhoto)
	return nil
}
