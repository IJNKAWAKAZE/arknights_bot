package system

import (
	"arknights_bot/config"
	"arknights_bot/utils/media"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"sync"
)

// fileId 缓存帮助图片，改并发执行后需要加锁保护。
var (
	fileId   string
	fileIdMu sync.Mutex
)

// HelpHandle 帮助
func HelpHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID

	_, _ = config.Arknights.SendChatAction(chatId, "upload_photo")

	fileIdMu.Lock()
	if fileId == "" {
		port := viper.GetString("http.port")
		pic, err := media.Screenshot("http://localhost:"+port+"/help", 0, 1.5)
		if err != nil {
			fileIdMu.Unlock()
			config.Arknights.ReplyText(chatId, messageId, err.Error())
			return nil
		}
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
		sendPhoto.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendPhoto)
		if err != nil {
			fileIdMu.Unlock()
			log.Println(err)
			return err
		}
		fileId = msg.Photo[0].FileID
		fileIdMu.Unlock()
		return nil
	}
	cachedFileId := fileId
	fileIdMu.Unlock()

	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(cachedFileId))
	sendPhoto.ReplyToMessageID = messageId
	config.Arknights.Send(sendPhoto)
	return nil
}
