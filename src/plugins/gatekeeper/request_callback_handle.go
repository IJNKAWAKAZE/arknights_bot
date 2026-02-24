package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"strconv"
	"strings"
)

func RequestCallBackData(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")

	if len(d) < 4 {
		return nil
	}

	userId, _ := strconv.ParseInt(d[1], 10, 64)
	chatId, _ := strconv.ParseInt(d[2], 10, 64)

	if has, correct := verifySet.checkExistAndRemove(userId, chatId); has {
		if d[3] != correct {
			callbackQuery.Answer(true, "验证未通过")
			bot.Arknights.DeclineChatJoinRequest(chatId, userId)
		} else {
			callbackQuery.Answer(true, "验证通过！")
			bot.Arknights.ApproveChatJoinRequest(chatId, userId)
			// 新人入群提醒
			var joined utils.GroupJoined
			utils.GetJoinedByChatId(chatId).Scan(&joined)
			var welcome string
			if joined.Welcome != "" {
				welcome = "，" + joined.Welcome
			}
			text := fmt.Sprintf("欢迎[%s](tg://user?id=%d)%s\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, callbackQuery.From.FullName()), callbackQuery.From.ID, welcome)
			if joined.Reg != -1 {
				chat, _ := bot.Arknights.GetChatInfo(chatId)
				if chat.UserName != "" {
					text += fmt.Sprintf("建议阅读群公约：[点击阅读](https://t.me/%s/%d)", chat.UserName, joined.Reg)
				} else {
					text += fmt.Sprintf("建议阅读群公约：[点击阅读](https://t.me/c/%s/%d)", strings.ReplaceAll(strconv.FormatInt(chat.ID, 10), "-100", ""), joined.Reg)
				}
			}
			sendMessage := tgbotapi.NewMessage(chatId, text)
			sendMessage.ParseMode = tgbotapi.ModeMarkdownV2
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(chatId, msg.MessageID, 3600)
		}
		callbackQuery.Delete()
	}
	return nil
}
