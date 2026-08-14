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
		config.Arknights.ReplyText(chatId, messageId, "国际服暂不可用")
		return nil
	}
	_, _ = config.Arknights.SendChatAction(chatId, "upload_photo")

	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/card?userId=%d&uid=%s&sklandId=%s", port, userAccount.UserNumber, uid, userAccount.SklandId), 0, 1)
	if err != nil {
		config.Arknights.SendMarkdownV2(chatId, err.Error(), messageId)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	/*sendPhoto.Caption = "点击复制UID:`" + uid + "`"
	sendPhoto.ParseMode = tgbotapi.ModeMarkdownV2*/
	config.Arknights.Send(sendPhoto)
	return nil
}
