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

type B struct {
	Size   int    `json:"size"`
	FileId string `json:"fileId"`
}

var BoxMap = make(map[int64]B)

type PlayerOperationBox struct {
	commandOperation.OperationAbstract
}

// BoxHandle 我的干员

func (_ PlayerOperationBox) Run(uid string, userAccount account.UserAccount, chatId int64, message *tgbotapi.Message) (bool, error) {
	messageId := message.MessageID
	param := message.CommandArguments()
	sendAction := tgbotapi.NewChatAction(chatId, "upload_document")
	bot.Arknights.Send(sendAction)

	matched, _ := regexp.MatchString("^[0-9\\d]+(,[0-9\\d]+)*$", param)
	if param != "" && param != "all" && !matched {
		sendMessage := tgbotapi.NewMessage(chatId, "参数错误")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/box?userId=%d&uid=%s&param=%s", port, userAccount.UserNumber, uid, param), 0)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败，token可能已失效请重设token。")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	// BOX有改变
	if BoxMap[userAccount.UserNumber].Size != len(pic) {
		sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "box.png"})
		sendDocument.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendDocument)
		if err == nil {
			b := B{
				Size:   len(pic),
				FileId: msg.Document.FileID,
			}
			BoxMap[userAccount.UserNumber] = b
		}
		return true, nil
	}
	// BOX无改变
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "box.png"})
	if BoxMap[userAccount.UserNumber].FileId != "" {
		sendDocument.BaseFile = tgbotapi.BaseFile{
			BaseChat: tgbotapi.BaseChat{
				ChatID: chatId,
			},
			File: tgbotapi.FileID(BoxMap[userAccount.UserNumber].FileId),
		}
	}
	sendDocument.ReplyToMessageID = messageId
	msg, _ := bot.Arknights.Send(sendDocument)
	b := B{
		Size:   len(pic),
		FileId: msg.Document.FileID,
	}
	BoxMap[userAccount.UserNumber] = b
	return true, nil
}
