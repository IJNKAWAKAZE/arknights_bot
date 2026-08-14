package gatekeeper

import (
	"arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
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
			auditJoin(chatId, userId, callbackQuery.From.FullName(), "拒绝", "入群申请验证答错")
			config.Arknights.DeclineChatJoinRequest(chatId, userId)
			log.Printf("入群验证：拒绝用户 %d（%s）加入群 %d，原因：答错", userId, callbackQuery.From.FullName(), chatId)
		} else {
			callbackQuery.Answer(true, "验证通过！")
			auditJoin(chatId, userId, callbackQuery.From.FullName(), "验证通过", "入群申请验证答对")
			config.Arknights.ApproveChatJoinRequest(chatId, userId)
			log.Printf("入群验证：通过用户 %d（%s）加入群 %d，原因：答对", userId, callbackQuery.From.FullName(), chatId)
			// 新人入群提醒
			var joined model.GroupJoined
			repo.GetJoinedByChatId(chatId).Scan(&joined)
			var welcome string
			if joined.Welcome != "" {
				welcome = "，" + joined.Welcome
			}
			text := fmt.Sprintf("欢迎[%s](tg://user?id=%d)%s\n", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, callbackQuery.From.FullName()), callbackQuery.From.ID, welcome)
			if joined.Reg != -1 {
				chat, _ := config.Arknights.GetChatInfo(chatId)
				if chat.UserName != "" {
					text += fmt.Sprintf("建议阅读群公约：[点击阅读](https://t.me/%s/%d)", chat.UserName, joined.Reg)
				} else {
					text += fmt.Sprintf("建议阅读群公约：[点击阅读](https://t.me/c/%s/%d)", strings.ReplaceAll(strconv.FormatInt(chat.ID, 10), "-100", ""), joined.Reg)
				}
			}
			msg, err := config.Arknights.SendMarkdownV2(chatId, text)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(chatId, msg.MessageID, 3600)
		}
		callbackQuery.Delete()
	}
	return nil
}
