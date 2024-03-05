package operator

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
)

// OperatorHandle 干员查询
func OperatorHandle(update tgbotapi.Update) (bool, error) {
	text := "干员-"
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	name := update.Message.CommandArguments()
	if name == "" {
		delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
		bot.Arknights.Send(delMsg)
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
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return true, nil
	}
	operator := ParseOperator(name)
	if operator.OP.Name == "" {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "查无此人，请输入正确的干员名称。")
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, bot.MsgDelDelay)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

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
	if utils.RedisIsExists(key) {
		fileId = utils.RedisGet(key)
	}

	if fileId != "" {
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(fileId))
		sendPhoto.ReplyToMessageID = messageId
		sendPhoto.ReplyMarkup = inlineKeyboardMarkup
		bot.Arknights.Send(sendPhoto)
		return true, nil
	}

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/operator?name=%s", port, name), 0, 1.5)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败，请重试。")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return true, nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyMarkup = inlineKeyboardMarkup
	sendPhoto.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		log.Println(err)
		return true, err
	}
	utils.RedisSet(key, msg.Photo[0].FileID, 0)
	return true, nil
}
