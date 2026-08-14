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

type PlayerOperationState struct {
	commandoperation.OperationAbstract
}

func (_ PlayerOperationState) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) error {
	messageId := message.MessageID
	if userAccount.ServerName == "国际服" {
		config.Arknights.ReplyText(chatId, messageId, "国际服暂不可用")
		return nil
	}
	_, _ = config.Arknights.SendChatAction(chatId, "upload_photo")

	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/state?userId=%d&uid=%s&sklandId=%s", port, userAccount.UserNumber, uid, userAccount.SklandId), 0, 1)
	if err != nil {
		config.Arknights.SendMarkdownV2(chatId, err.Error(), messageId)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	config.Arknights.Send(sendPhoto)
	return nil
}
