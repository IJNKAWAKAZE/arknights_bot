package gatekeeper

import (
	"arknights_bot/config"
	"arknights_bot/utils/media"
	"arknights_bot/utils/model"
	"arknights_bot/utils/search"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"math/big"
	"time"
)

var verifyC = make(chan interface{}, 10)

func VerifyMember(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	userId := message.From.ID
	name := message.From.FullName()
	messageId := message.MessageID
	// 限制用户发送消息
	_, err := config.Arknights.RestrictChatMember(chatId, userId, tgbotapi.NoMessagesPermission)
	if err != nil {
		log.Println(err.Error())
		return
	}

	// 抽取验证信息
	operatorsPool := search.GetOperators()
	if len(operatorsPool) < 12 { // 不足 12 个干员时去重循环会无限空转，恢复权限后放弃本次验证
		log.Println("入群验证：干员数据不足，恢复用户", userId, "发言权限")
		auditJoin(chatId, userId, name, "未验证放行", "干员数据不足，恢复发言权限")
		config.Arknights.RestrictChatMember(chatId, userId, tgbotapi.AllPermissions)
		return
	}
	var randNumMap = make(map[int64]struct{})
	var options []model.Operator
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
		if len(operator.Skins) == 0 {
			i--
			continue
		}
		painting := operator.Skins[0].Url
		if painting != "" {
			options = append(options, model.Operator{
				Name:     operatorName,
				ThumbURL: painting,
			})
		} else {
			i--
		}
	}

	if len(options) < 2 {
		log.Println("入群验证：可用干员不足，恢复用户", userId, "发言权限")
		auditJoin(chatId, userId, name, "未验证放行", "可用干员不足，恢复发言权限")
		verifySet.checkExistAndRemove(userId, chatId)
		config.Arknights.RestrictChatMember(chatId, userId, tgbotapi.AllPermissions)
		return
	}
	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options)-1)))
	correct := options[r.Int64()+1]
	verifySet.add(userId, chatId, correct.Name)

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
	if len(verifyC) > 0 {
		obj := <-verifyC
		log.Println(obj, "停止发送验证信息")
		message.Delete()
		config.Arknights.BanChatMember(chatId, userId)
		verifySet.checkExistAndRemove(userId, chatId)
		return
	}
	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: media.GetImg(correct.ThumbURL)})
	sendPhoto.ReplyMarkup = inlineKeyboardMarkup
	sendPhoto.Caption = fmt.Sprintf("欢迎[%s](tg://user?id=%d)，请选择上图干员的正确名字，60秒未选择自动踢出。", tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, name), userId)
	sendPhoto.ParseMode = tgbotapi.ModeMarkdownV2
	photo, err := config.Arknights.Send(sendPhoto)
	if err != nil {
		log.Printf("发送图片失败：%s，原因：%s", correct.ThumbURL, err.Error())
		auditJoin(chatId, userId, name, "未验证放行", "发送验证图片失败，恢复发言权限")
		config.Arknights.RestrictChatMember(chatId, userId, tgbotapi.AllPermissions)
		verifySet.checkExistAndRemove(userId, chatId)
		return
	}
	go verify(chatId, userId, photo.MessageID, messageId)
}

func unban(chatId, userId int64) {
	time.Sleep(time.Minute)
	config.Arknights.UnbanChatMember(chatId, userId)
}

func verify(chatId int64, userId int64, messageId int, joinMessageId int) {
	time.Sleep(time.Minute)
	if has, _ := verifySet.checkExistAndRemove(userId, chatId); !has {
		return
	}

	// 踢出超时未验证用户
	config.Arknights.BanChatMember(chatId, userId)
	// 删除用户入群提醒
	delJoinMessage := tgbotapi.NewDeleteMessage(chatId, joinMessageId)
	config.Arknights.Send(delJoinMessage)
	// 删除入群验证消息
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	config.Arknights.Send(delMsg)
	time.Sleep(time.Minute)
	// 解除用户封禁
	config.Arknights.UnbanChatMember(chatId, userId)
}
