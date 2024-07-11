package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"math/big"
	"time"
)

func VerifyMember(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	userId := message.From.ID
	name := message.From.FullName()
	messageId := message.MessageID
	// 限制用户发送消息
	_, err := bot.Arknights.RestrictChatMember(chatId, userId, tgbotapi.NoMessagesPermission)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 抽取验证信息
	operatorsPool := utils.GetOperators()
	var randNumMap = make(map[int64]struct{})
	var options []utils.Operator
	for i := 0; i < 12; i++ { // 随机抽取 12 个干员
		var operatorIndex int64
		for { // 抽到重复索引则重新抽取
			r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(operatorsPool))))
			if _, has := randNumMap[r.Int64()]; !has {
				operatorIndex = r.Int64()
				randNumMap[operatorIndex] = struct{}{}
				break
			}
		}
		operator := operatorsPool[operatorIndex]
		operatorName := operator.Name
		painting := operator.Skins[0].Url
		if painting != "" {
			options = append(options, utils.Operator{
				Name:     operatorName,
				ThumbURL: painting,
			})
		} else {
			i--
		}
	}

	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options)-1)))
	correct := options[r.Int64()+1]

	var buttons [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(options); i += 2 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(options[i].Name, fmt.Sprintf("verify,%d,%s,%d", userId, options[i].Name, messageId)),
			tgbotapi.NewInlineKeyboardButtonData(options[i+1].Name, fmt.Sprintf("verify,%d,%s,%d", userId, options[i+1].Name, messageId)),
		))
	}
	buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("✅放行", fmt.Sprintf("verify,%d,PASS,%d", userId, messageId)),
		tgbotapi.NewInlineKeyboardButtonData("🚫封禁", fmt.Sprintf("verify,%d,BAN,%d", userId, messageId)),
	))
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: utils.GetImg(correct.ThumbURL)})
	sendPhoto.ReplyMarkup = inlineKeyboardMarkup
	sendPhoto.Caption = fmt.Sprintf("欢迎[%s](tg://user?id=%d)，请选择上图干员的正确名字，60秒未选择自动踢出。", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, name), userId)
	sendPhoto.ParseMode = tgbotapi.ModeMarkdownV2
	photo, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		log.Printf("发送图片失败：%s，原因：%s", correct.ThumbURL, err.Error())
		bot.Arknights.RestrictChatMember(chatId, userId, tgbotapi.AllPermissions)
		return
	}
	verifySet.add(userId, chatId, correct.Name)
	go verify(chatId, userId, photo.MessageID, messageId)
}

func unban(chatId, userId int64) {
	time.Sleep(time.Minute)
	if bot.Arknights.GetChatMemberStatus(chatId, userId) != tgbotapi.KICKED {
		bot.Arknights.UnbanChatMember(chatId, userId)
	}
}

func verify(chatId int64, userId int64, messageId int, joinMessageId int) {
	time.Sleep(time.Minute)
	if has, _ := verifySet.checkExistAndRemove(userId, chatId); !has {
		return
	}

	// 踢出超时未验证用户
	bot.Arknights.BanChatMember(chatId, userId)
	// 删除用户入群提醒
	delJoinMessage := tgbotapi.NewDeleteMessage(chatId, joinMessageId)
	bot.Arknights.Send(delJoinMessage)
	// 删除入群验证消息
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)
	time.Sleep(time.Minute)
	// 解除用户封禁
	if bot.Arknights.GetChatMemberStatus(chatId, userId) != tgbotapi.KICKED {
		bot.Arknights.UnbanChatMember(chatId, userId)
	}
}
