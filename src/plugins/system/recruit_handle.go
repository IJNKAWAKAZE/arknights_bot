package system

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"strings"
)

func RecruitHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	caption := strings.Split(update.Message.Caption, " ")
	param := caption[len(caption)-1]
	photos := update.Message.Photo
	return recruit(chatId, messageId, param, photos[len(photos)-1].FileID)
}

func ReplyRecruitHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	text := strings.Split(update.Message.Text, " ")
	param := text[len(text)-1]
	photos := update.Message.ReplyToMessage.Photo
	return recruit(chatId, messageId, param, photos[len(photos)-1].FileID)
}

func recruit(chatId int64, messageId int, param, fileId string) error {
	var tags []string
	file, _ := utils.DownloadFile(fileId)
	lang, engine, sep := "chs", "2", "\n"
	if param == "jp" {
		lang, engine, sep = "jpn", "1", "\r\n"
	}

	results, err := utils.OCR(file, lang, engine, sep)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, "识别失败请稍后再试")
		bot.Arknights.Send(sendMessage)
	}
	if results == nil {
		log.Println("图片识别失败")
		return nil
	}
	for _, result := range results {
		if bot.RecruitTagMap[result] != "" {
			log.Println(bot.RecruitTagMap[result])
			tags = append(tags, bot.RecruitTagMap[result])
		}
	}
	if len(tags) != 5 {
		sendMessage := tgbotapi.NewMessage(chatId, "标签数量错误，请更换图片。")
		sendMessage.ReplyToMessageID = messageId
		_, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		return nil
	}
	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	port := viper.GetString("http.port")
	pic, err := utils.Screenshot(fmt.Sprintf("http://localhost:%s/recruit?tags=%s&client=%s", port, strings.Join(tags, " "), param), 0, 1.5)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ReplyToMessageID = messageId
		_, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		return nil
	}
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "recruit.jpg"})
	sendDocument.ReplyToMessageID = messageId
	_, err = bot.Arknights.Send(sendDocument)
	if err != nil {
		return err
	}
	return nil
}
