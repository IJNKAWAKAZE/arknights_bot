package skin

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/tidwall/sjson"
)

// SkinHandle 干员皮肤查询
func SkinHandle(update tgbotapi.Update) error {
	text := "皮肤-"
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
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}
	operator := utils.GetOperatorByName(name)
	if operator.Name == "" {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "查无此人，请输入正确的干员名称。")
		sendMessage.ReplyToMessageID = messageId
		msg, err := bot.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, bot.MsgDelDelay)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	bot.Arknights.Send(sendAction)

	content := "[]"
	for _, skin := range operator.Skins {
		src, _ := sjson.Set("", "attrs.src", skin.Url)
		attrs, _ := sjson.Set(src, "tag", "img")
		content, _ = sjson.SetRaw(content, "-1", attrs)
	}
	skinUrl := utils.CreateTelegraphPage(content, name+"的皮肤")
	sendMessage := tgbotapi.NewMessage(chatId, skinUrl)
	sendMessage.ReplyToMessageID = messageId
	bot.Arknights.Send(sendMessage)
	return nil
}
