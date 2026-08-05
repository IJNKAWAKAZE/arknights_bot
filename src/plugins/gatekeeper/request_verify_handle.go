package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"math/big"
	"strconv"
	"time"
)

func VerifyRequestMember(update tgbotapi.Update) {
	chatId := update.ChatJoinRequest.Chat.ID
	userId := update.ChatJoinRequest.From.ID
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

	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options)-1)))
	correctIdx := r.Int64() + 1
	correct := options[correctIdx]
	// 原子登记验证条目（防止重投更新产生重复题目），callback 携带选项序号而非干员名，
	// 既规避 Telegram 64 字节回调数据上限，也避免答案明文出现在客户端可见数据中
	if !verifySet.addIfNotExist(userId, chatId, strconv.FormatInt(correctIdx, 10)) {
		return
	}

	var buttons [][]tgbotapi.InlineKeyboardButton
	for i := 0; i < len(options); i += 2 {
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(options[i].Name, fmt.Sprintf("request,%d,%d,%d", userId, chatId, i)),
			tgbotapi.NewInlineKeyboardButtonData(options[i+1].Name, fmt.Sprintf("request,%d,%d,%d", userId, chatId, i+1)),
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
