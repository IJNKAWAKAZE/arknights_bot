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

func VerifyRequestMember(update tgbotapi.Update) {
	chatId := update.ChatJoinRequest.Chat.ID
	userId := update.ChatJoinRequest.From.ID
	member := update.ChatJoinRequest.From
	// 广告词检查：昵称/用户名/简介命中广告词直接拒绝
	bio := ""
	if chat, err := bot.Arknights.GetChatInfo(userId); err == nil {
		bio = chat.Bio
	}
	if isAdUser(member, bio) {
		log.Printf("入群验证：用户 %d（%s）命中广告词，拒绝入群申请。群：%d", userId, member.FullName(), chatId)
		bot.Arknights.DeclineChatJoinRequest(chatId, userId)
		return
	}
	if verifySet.checkExist(userId, chatId) {
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

	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
	correct := options[r.Int64()]
	verifySet.add(userId, chatId, correct.Name)

	var buttons [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(options); i += 2 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(options[i].Name, fmt.Sprintf("request,%d,%d,%s", userId, chatId, options[i].Name)),
			tgbotapi.NewInlineKeyboardButtonData(options[i+1].Name, fmt.Sprintf("request,%d,%d,%s", userId, chatId, options[i+1].Name)),
		))
	}
	inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
		buttons...,
	)
	sendPhoto := tgbotapi.NewPhoto(userId, tgbotapi.FileBytes{Bytes: utils.GetImg(correct.ThumbURL)})
	sendPhoto.ReplyMarkup = inlineKeyboardMarkup
	sendPhoto.Caption = "请选择上图干员的正确名字"
	photo, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		log.Printf("发送图片失败：%s，原因：%s", correct.ThumbURL, err.Error())
		bot.Arknights.ApproveChatJoinRequest(chatId, userId)
		verifySet.checkExistAndRemove(userId, chatId)
		return
	}
	go requestVerify(chatId, userId, photo.MessageID)
}

func requestVerify(chatId int64, userId int64, messageId int) {
	time.Sleep(time.Minute)
	if has, _ := verifySet.checkExistAndRemove(userId, chatId); !has {
		return
	}
	bot.Arknights.DeclineChatJoinRequest(chatId, userId)
	// 删除入群验证消息
	delMsg := tgbotapi.NewDeleteMessage(userId, messageId)
	bot.Arknights.Send(delMsg)
}
