package player

import (
	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils/media"
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
	if userAccount.ServerName == "国际服" {
		sendMessage := tgbotapi.NewMessage(chatId, "国际服暂不可用")
		sendMessage.ReplyToMessageID = messageId
		config.Arknights.Send(sendMessage)
		return nil
	}
	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	config.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/card?userId=%d&uid=%s&sklandId=%s", port, userAccount.UserNumber, uid, userAccount.SklandId), 0, 1)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		sendMessage.ReplyToMessageID = messageId
		config.Arknights.Send(sendMessage)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	/*sendPhoto.Caption = "点击复制UID:`" + uid + "`"
	sendPhoto.ParseMode = tgbotapi.ModeMarkdownV2*/
	config.Arknights.Send(sendPhoto)
	return nil
}
