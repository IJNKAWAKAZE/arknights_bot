package material

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"log"
)

// MaterialHandle 材料查询
func MaterialHandle(update tgbotapi.Update) error {
	text := "材料-"
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	name := update.Message.CommandArguments()
	if name == "" {
		delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
		bot.Arknights.Send(delMsg)
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:                         "选择材料",
					SwitchInlineQueryCurrentChat: &text,
				},
			),
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要查询的材料")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	materials := utils.GetItemByName(name)
	if len(materials) == 0 {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "未查询到此材料，请输入正确的材料名称。")
		sendMessage.ReplyToMessageID = messageId
		msg, _ := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, bot.MsgDelDelay)
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	fileId := ""
	key := "material:" + name
	if utils.RedisIsExists(key) {
		fileId = utils.RedisGet(key)
	}

	if fileId != "" {
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileID(fileId))
		sendPhoto.ReplyToMessageID = messageId
		bot.Arknights.Send(sendPhoto)
		return nil
	}

	port := viper.GetString("http.port")
	pic := utils.Screenshot(fmt.Sprintf("http://localhost:%s/material?name=%s", port, name), 0, 1.5)
	if pic == nil {
		sendMessage := tgbotapi.NewMessage(chatId, "生成图片失败，请重试。")
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	sendPhoto.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		log.Println(err)
		return err
	}
	utils.RedisSet(key, msg.Photo[0].FileID, 0)
	return nil
}
