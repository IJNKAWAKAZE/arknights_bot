package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// ReportHandle 举报
func ReportHandle(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID
	messageId := message.MessageID

	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)

	if message.ReplyToMessage != nil {
		replyToMessage := message.ReplyToMessage
		replyMessageId := replyToMessage.MessageID
		target := replyToMessage.From.ID
		name := utils.GetFullName(replyToMessage.From)

		if utils.IsAdmin(chatId, target) {
			sendMessage := tgbotapi.NewMessage(chatId, "无法举报管理员！")
			sendMessage.ReplyToMessageID = messageId
			msg, _ := bot.Arknights.Send(sendMessage)
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}

		// 获取全部管理员
		getAdmins := tgbotapi.ChatAdministratorsConfig{
			ChatConfig: tgbotapi.ChatConfig{
				ChatID: chatId,
			},
		}

		var buttons [][]tgbotapi.InlineKeyboardButton

		var text bytes.Buffer
		text.WriteString(fmt.Sprintf("被举报人：[%s](tg://user?id=%d)\n", utils.EscapesMarkdownV2(name), target))
		text.WriteString(fmt.Sprintf("消息存放：[%d](https://t.me/%s/%d)", replyMessageId, replyToMessage.Chat.UserName, replyMessageId))
		charAdmins, _ := bot.Arknights.GetChatAdministrators(getAdmins)
		var admins []int64
		for _, admin := range charAdmins {
			if !admin.User.IsBot {
				admins = append(admins, admin.User.ID)
			}
		}

		for _, admin := range admins {
			text.WriteString(fmt.Sprintf("[\u200b](tg://user?id=%d) ", admin))
		}

		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("🚫封禁", fmt.Sprintf("report,%s,%d", "BAN", target)),
			tgbotapi.NewInlineKeyboardButtonData("❌关闭", fmt.Sprintf("report,%s,%d", "CLOSE", target)),
		))

		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)

		sendMessage := tgbotapi.NewMessage(chatId, text.String())
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		bot.Arknights.Send(sendMessage)
	}

	return nil
}
