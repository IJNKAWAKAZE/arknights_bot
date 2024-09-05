package enemy

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/spf13/viper"
	"log"
	"net/url"
)

// EnemyHandle 敌人查询
func EnemyHandle(update tgbotapi.Update) error {
	text := "敌人-"
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID
	name := update.Message.CommandArguments()
	if name == "" {
		update.Message.Delete()
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.InlineKeyboardButton{
					Text:                         "选择敌人",
					SwitchInlineQueryCurrentChat: &text,
				},
			),
		)
		sendMessage := tgbotapi.NewMessage(chatId, "请选择要查询的敌人")
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	if name == "多萝西" {
		name = "多萝西(敌方)"
	}
	enemy := ParseEnemy(name)
	if enemy.Name == "" {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "未查询到此敌人，请输入正确的敌人名称。")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, bot.MsgDelDelay)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	bot.Arknights.Send(sendAction)

	link := viper.GetString("api.wiki") + url.PathEscape(name)
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.InlineKeyboardButton{
				Text: "查看详情",
				URL:  &link,
			},
		),
	)

	fileId := ""
	key := "enemy:" + name
	if utils.RedisIsExists(key) {
		fileId = utils.RedisGet(key)
	}

	if fileId != "" {
		sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileID(fileId))
		sendDocument.ReplyToMessageID = messageId
		sendDocument.ReplyMarkup = inlineKeyboardMarkup
		bot.Arknights.Send(sendDocument)
		return nil
	}

	port := viper.GetString("http.port")
	pic, err := utils.Screenshot(fmt.Sprintf("http://localhost:%s/enemy?name=%s", port, name), 0, 1.5)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ReplyToMessageID = messageId
		bot.Arknights.Send(sendMessage)
		return nil
	}
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "enemy.jpg"})
	sendDocument.ReplyMarkup = inlineKeyboardMarkup
	sendDocument.ReplyToMessageID = messageId
	msg, err := bot.Arknights.Send(sendDocument)
	if err != nil {
		log.Println(err)
		return err
	}
	utils.RedisSet(key, msg.Document.FileID, 0)
	return nil
}
