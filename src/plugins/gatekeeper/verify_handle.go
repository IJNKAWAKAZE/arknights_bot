package gatekeeper

import (
	bot "arknights_bot/config"
	"arknights_bot/utils"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/big"
	"strconv"
	"time"
)

func VerifyMember(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	userId := message.From.ID
	name := utils.GetFullName(message.From)
	for _, m := range message.NewChatMembers {

		// 限制用户发送消息
		restrictChatMemberConfig := tgbotapi.RestrictChatMemberConfig{
			Permissions: &tgbotapi.ChatPermissions{
				CanSendMessages: false,
			},
			ChatMemberConfig: tgbotapi.ChatMemberConfig{
				ChatID: chatId,
				UserID: m.ID,
			},
		}
		_, err := bot.Arknights.Request(restrictChatMemberConfig)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// 抽取验证信息
		operatorsPool := utils.GetOperators()
		var operatorMap = make(map[string]struct{})
		var randNumMap = make(map[int64]struct{})
		var options []Verify
		for i := 0; i < 4; i++ { // 随机抽取 4 个干员
			var operatorIndex int64
			for { // 抽到重复索引则重新抽取
				r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(operatorsPool))))
				if _, has := randNumMap[r.Int64()]; !has {
					operatorIndex = r.Int64()
					break
				}
			}
			operator := operatorsPool[operatorIndex]
			shipName := operator.Get("name").String()
			painting := operator.Get("painting").String()
			if painting != "" {
				if _, has := operatorMap[shipName]; has { // 如果 map 中已存在该干员，则跳过
					continue
				}
				// 保存干员信息
				operatorMap[shipName] = struct{}{}
				options = append(options, Verify{
					Name:     shipName,
					Painting: painting,
				})
			}
		}

		r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
		random, _ := strconv.Atoi(r.String())
		correct := options[random]

		var buttons [][]tgbotapi.InlineKeyboardButton
		userIdStr := strconv.FormatInt(userId, 10)
		for _, v := range options {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(v.Name, userIdStr+","+v.Name+","+correct.Name),
			))
		}
		buttons = append(buttons, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅放行", userIdStr+",PASS,"+name),
			tgbotapi.NewInlineKeyboardButtonData("🚫封禁", userIdStr+",BAN,"+name),
		))
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(correct.Painting))
		sendPhoto.ReplyMarkup = inlineKeyboardMarkup
		sendPhoto.Caption = "欢迎<a href=\"tg://user?id=" + userIdStr + "\">" + name + "</a>，请选择上图干员的正确名字，60秒未选择自动踢出。"
		sendPhoto.ParseMode = tgbotapi.ModeHTML
		photo, err := bot.Arknights.Send(sendPhoto)
		if err != nil {
			log.Println(err)
			restrictChatMemberConfig = tgbotapi.RestrictChatMemberConfig{
				Permissions: &tgbotapi.ChatPermissions{
					CanSendMessages:       true,
					CanSendMediaMessages:  true,
					CanSendPolls:          true,
					CanSendOtherMessages:  true,
					CanAddWebPagePreviews: true,
					CanInviteUsers:        true,
					CanChangeInfo:         true,
					CanPinMessages:        true,
				},
				ChatMemberConfig: tgbotapi.ChatMemberConfig{
					ChatID: chatId,
					UserID: userId,
				},
			}
			bot.Arknights.Send(restrictChatMemberConfig)
			return
		}
		val := fmt.Sprintf("verify%d%d", chatId, userId)
		utils.RedisAddSet("verify", val)
		go verify(val, chatId, userId, photo.MessageID, name)
	}
}

func unban(chatMember tgbotapi.ChatMemberConfig) {
	time.Sleep(time.Minute)
	unbanChatMemberConfig := tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: chatMember,
		OnlyIfBanned:     true,
	}
	bot.Arknights.Send(unbanChatMemberConfig)
}

func verify(val string, chatId int64, userId int64, messageId int, name string) {
	time.Sleep(time.Minute)
	if !utils.RedisSetIsExists("verify", val) {
		return
	}
	chatMember := tgbotapi.ChatMemberConfig{ChatID: chatId, UserID: userId}
	kickChatMemberConfig := tgbotapi.KickChatMemberConfig{
		ChatMemberConfig: chatMember,
	}
	bot.Arknights.Send(kickChatMemberConfig)
	sendMessage := tgbotapi.NewMessage(chatId, fmt.Sprintf("<a href=\"tg://user?id=%d\">%s</a>超时未验证，已被踢出。", userId, name))
	sendMessage.ParseMode = tgbotapi.ModeHTML
	msg, _ := bot.Arknights.Send(sendMessage)
	utils.AddDelQueue(msg.Chat.ID, msg.MessageID, 1)
	utils.RedisDelSetItem("verify", val)
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)
	time.Sleep(time.Minute)
	bot.Arknights.Send(tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: chatMember,
		OnlyIfBanned:     true,
	})
}
