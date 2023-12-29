package handle

import (
	"arknights_bot/bot/modules"
	"arknights_bot/bot/utils"
	"crypto/rand"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/big"
	"strconv"
	"time"
)

func Verify(message *tgbotapi.Message) {
	chatId := message.Chat.ID
	userId := message.From.ID
	name := utils.GetFullName(message.From)
	var operators = make(map[string]modules.Verify)
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
		_, err := utils.SetMemberPermissions(restrictChatMemberConfig)
		if err != nil {
			log.Println(err.Error())
			return
		}

		// 抽取验证信息
		operatorsPool := utils.GetOperators()
		var operatorMap = make(map[string]struct{})
		var randNumMap = make(map[int64]struct{})
		var options []modules.Verify
		for i := 0; i < 4; i++ { // 随机抽取 4 个干员
			var operatorIndex int64
			for { // 抽到重复索引则重新抽取
				r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(operatorsPool))))
				if _, has := randNumMap[r.Int64()]; !has {
					operatorIndex = r.Int64()
					break
				}
			}
			ship := operatorsPool[operatorIndex]
			shipName := ship.Get("name").String()
			painting := ship.Get("painting").String()
			if painting != "" {
				if _, has := operatorMap[shipName]; has { // 如果 map 中已存在该干员，则跳过
					continue
				}
				// 保存干员信息
				operatorMap[shipName] = struct{}{}
				options = append(options, modules.Verify{
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
		for _, v := range operators {
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
		photo, err := utils.SendPhoto(sendPhoto)
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
			utils.SetMemberPermissions(restrictChatMemberConfig)
			return
		}
		val := "verify" + strconv.FormatInt(chatId, 10) + userIdStr
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
	utils.UnbanChatMember(unbanChatMemberConfig)
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
	utils.KickChatMember(kickChatMemberConfig)
	sendMessage := tgbotapi.NewMessage(chatId, "<a href=\"tg://user?id="+strconv.FormatInt(userId, 10)+"\">"+name+"</a>超时未验证，已被踢出。")
	sendMessage.ParseMode = tgbotapi.ModeHTML
	msg, _ := utils.SendMessage(sendMessage)
	utils.AddDelQueue(msg.Chat.ID, msg.MessageID, 1)
	utils.RedisDelSetItem("verify", val)
	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	utils.DeleteMessage(delMsg)
	time.Sleep(time.Minute)
	utils.UnbanChatMember(tgbotapi.UnbanChatMemberConfig{
		ChatMemberConfig: chatMember,
		OnlyIfBanned:     true,
	})
}
