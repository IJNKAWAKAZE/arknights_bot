package system

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"bytes"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"strconv"
	"strings"
)

// ReportHandle 举报
func ReportHandle(update tgbotapi.Update) error {
	message := update.Message
	chatId := message.Chat.ID

	message.Delete()

	if message.ReplyToMessage != nil {
		replyToMessage := message.ReplyToMessage
		replyMessageId := replyToMessage.MessageID
		target := replyToMessage.From.ID
		name := replyToMessage.From.FullName()

		if config.Arknights.IsAdmin(chatId, target) {
			msg, err := config.Arknights.SendText(chatId, "无法举报管理员！")
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, config.MsgDelDelay)
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
		text.WriteString(fmt.Sprintf("被举报人：[%s](tg://user?id=%d)\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, name), target))
		if replyToMessage.Chat.UserName != "" {
			text.WriteString(fmt.Sprintf("消息存放：[%d](https://t.me/%s/%d)", replyMessageId, replyToMessage.Chat.UserName, replyMessageId))
		} else {
			text.WriteString(fmt.Sprintf("消息存放：[%d](https://t.me/c/%s/%d)", replyMessageId, strings.ReplaceAll(strconv.FormatInt(replyToMessage.Chat.ID, 10), "-100", ""), replyMessageId))
		}
		charAdmins, _ := config.Arknights.GetChatAdministrators(getAdmins)
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
			tgbotapi.NewInlineKeyboardButtonData("🚫封禁", fmt.Sprintf("report,%s,%d,%d", "BAN", target, replyMessageId)),
			tgbotapi.NewInlineKeyboardButtonData("❌关闭", fmt.Sprintf("report,%s,%d,%d", "CLOSE", target, replyMessageId)),
		))

		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)

		sendMessage := tgbotapi.NewMessage(chatId, text.String())
		sendMessage.ReplyMarkup = inlineKeyboardMarkup
		sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
		config.Arknights.Send(sendMessage)
	}

	return nil
}
