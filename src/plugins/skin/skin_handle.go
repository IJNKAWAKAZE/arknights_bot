package skin

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/media"
	"arknights_bot/utils/search"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"github.com/tidwall/sjson"
)

// SkinHandle 干员皮肤查询
func SkinHandle(update tgbotapi.Update) error {
	text := "皮肤-"
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
	operator := search.GetOperatorByName(name)
	if operator.Name == "" {
		sendMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "查无此人，请输入正确的干员名称。")
		sendMessage.ReplyToMessageID = messageId
		msg, err := config.Arknights.Send(sendMessage)
		messagecleaner.AddDelQueue(chatId, messageId, config.MsgDelDelay)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	config.Arknights.Send(sendAction)

	content := "[]"
	for _, skin := range operator.Skins {
		if skin.Name != "" {
			children, _ := sjson.Set("{}", "tag", "h4")
			children, _ = sjson.Set(children, "children", []string{skin.Name})
			content, _ = sjson.SetRaw(content, "-1", children)
		}
		src, _ := sjson.Set("", "attrs.src", skin.Url)
		attrs, _ := sjson.Set(src, "tag", "img")
		content, _ = sjson.SetRaw(content, "-1", attrs)
	}
	skinUrl := media.CreateTelegraphPage(content, name+"的皮肤")
	sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("[%s的皮肤](%s)", name, skinUrl))
	sendMessage.ReplyToMessageID = messageId
	sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
	config.Arknights.Send(sendMessage)
	return nil
}
