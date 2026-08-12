package enemy

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/cache"
	"arknights_bot/utils/media"
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
		msg, err := config.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}
	config.DataMu.RLock()
	enemyName, has := config.EnemyName[name]
	config.DataMu.RUnlock()
	if has {
		name = enemyName
	}
	enemy := ParseEnemy(name)
	if enemy.Name == "" {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "未查询到此敌人，请输入正确的敌人名称。")
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, config.MsgDelDelay)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "upload_photo")
	config.Arknights.Send(sendAction)

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
	if cache.RedisIsExists(key) {
		fileId = cache.RedisGet(key)
	}

	if fileId != "" {
		sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileID(fileId))
		sendDocument.ReplyToMessageID = messageId
		sendDocument.ReplyMarkup = inlineKeyboardMarkup
		config.Arknights.Send(sendDocument)
		return nil
	}

	port := viper.GetString("http.port")
	pic, err := media.Screenshot(fmt.Sprintf("http://localhost:%s/enemy?name=%s", port, name), 0, 1.5)
	if err != nil {
		sendMessage := tgbotapi.NewMessage(chatId, err.Error())
		sendMessage.ReplyToMessageID = messageId
		config.Arknights.Send(sendMessage)
		return nil
	}
	sendDocument := tgbotapi.NewDocument(chatId, tgbotapi.FileBytes{Bytes: pic, Name: "enemy.jpg"})
	sendDocument.ReplyMarkup = inlineKeyboardMarkup
	sendDocument.ReplyToMessageID = messageId
	msg, err := config.Arknights.Send(sendDocument)
	if err != nil {
		log.Println(err)
		return err
	}
	cache.RedisSet(key, msg.Document.FileID, 0)
	return nil
}
