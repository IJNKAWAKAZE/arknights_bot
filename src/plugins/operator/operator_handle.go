package operator

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/cache"
	"arknights_bot/utils/media"
	"fmt"
	"log"

	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
)

// OperatorHandle 干员查询
func OperatorHandle(update tgbotapi.Update) error {
	text := "干员-"
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	name := update.Message.CommandArguments()
	if name == "" {
		update.Message.Delete()
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:                         "选择干员",
					SwitchInlineQueryCurrentChat: &text,
				},
			),
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要查询的干员")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	operator := ParseOperator(name)
	if operator.OP.Name == "" {
		msg, err := config.Arknights.ReplyText(update.Message.Chat.ID, messageId, "查无此人，请输入正确的干员名称。")
		messagecleaner.AddDelQueue(chatId, messageId, config.MsgDelDelay)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	name = operator.OP.Name
	_, _ = config.Arknights.SendChatAction(chatId, "upload_photo")

	url := viper.GetString("api.wiki") + name
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: "查看详情",
				URL:  &url,
			},
		),
	)

	fileId := ""
	key := "operator:" + name
	if cache.RedisIsExists(key) {
		fileId = cache.RedisGet(key)
	}

	if fileId != "" {
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(fileId))
		sendPhoto.ReplyToMessageID = messageId
		sendPhoto.ReplyMarkup = inlineKeyboardMarkup
		config.Arknights.Send(sendPhoto)
		return nil
	}

	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/operator?name=%s", port, name), 0, 1.5)
	if err != nil {
		config.Arknights.ReplyText(chatId, messageId, err.Error())
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyMarkup = inlineKeyboardMarkup
	sendPhoto.ReplyToMessageID = messageId
	msg, err := config.Arknights.Send(sendPhoto)
	if err != nil {
		log.Println(err)
		return err
	}
	cache.RedisSet(key, msg.Photo[0].FileID, 0)
	return nil
}
