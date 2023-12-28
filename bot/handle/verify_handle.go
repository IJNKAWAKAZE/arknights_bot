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
	var correct modules.Verify
	var options []modules.Verify
	var buttons [][]tgbotapi.InlineKeyboardButton
	for _, m := range message.NewChatMembers {
		chatPermissions := tgbotapi.ChatPermissions{
			CanSendMessages: false,
		}
		restrictChatMemberConfig := tgbotapi.RestrictChatMemberConfig{
			Permissions: &chatPermissions,
		}
		restrictChatMemberConfig.ChatID = chatId
		restrictChatMemberConfig.UserID = m.ID
		utils.SetMemberPermissions(restrictChatMemberConfig)
		operatorsPool := utils.GetOperators()
		for true {
			if len(operators) == 4 {
				break
			}
			r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(operatorsPool))))
			random, _ := strconv.Atoi(r.String())
			ship := operatorsPool[random]
			name := ship.Get("name").String()
			painting := ship.Get("painting").String()
			if painting != "" {
				var s = modules.Verify{
					Name:     name,
					Painting: painting,
				}
				operators[name] = s
			}
		}
		for _, v := range operators {
			options = append(options, v)
		}
		r, _ := rand.Int(rand.Reader, big.NewInt(4))
		random, _ := strconv.Atoi(r.String())
		correct = options[random]
		for _, v := range operators {
			btn := tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData(v.Name, strconv.FormatInt(userId, 10)+","+v.Name+","+correct.Name),
			)
			buttons = append(buttons, btn)
		}
		adminBtn := tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("✅放行", strconv.FormatInt(userId, 10)+",PASS,"+name),
			tgbotapi.NewInlineKeyboardButtonData("🚫封禁", strconv.FormatInt(userId, 10)+",BAN,"+name),
		)
		buttons = append(buttons, adminBtn)
		inlineKeyboardMarkup := tgbotapi.NewInlineKeyboardMarkup(
			buttons...,
		)
		sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(correct.Painting))
		sendPhoto.ReplyMarkup = inlineKeyboardMarkup
		sendPhoto.Caption = "欢迎<a href=\"tg://user?id=" + strconv.FormatInt(userId, 10) + "\">" + name + "</a>，请选择上图干员的正确名字，60秒未选择自动踢出。"
		sendPhoto.ParseMode = tgbotapi.ModeHTML
		photo, err := utils.SendPhoto(sendPhoto)
		if err != nil {
			log.Println(err)
			chatPermissions = tgbotapi.ChatPermissions{
				CanSendMessages:       true,
				CanSendMediaMessages:  true,
				CanSendPolls:          true,
				CanSendOtherMessages:  true,
				CanAddWebPagePreviews: true,
				CanInviteUsers:        true,
				CanChangeInfo:         true,
				CanPinMessages:        true,
			}
			restrictChatMemberConfig = tgbotapi.RestrictChatMemberConfig{
				Permissions: &chatPermissions,
			}
			restrictChatMemberConfig.ChatID = chatId
			restrictChatMemberConfig.UserID = userId
			utils.SetMemberPermissions(restrictChatMemberConfig)
			return
		}
		cid := strconv.FormatInt(chatId, 10)
		uid := strconv.FormatInt(userId, 10)
		val := "verify" + cid + uid
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
	if utils.RedisSetIsExists("verify", val) {
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
		unbanChatMemberConfig := tgbotapi.UnbanChatMemberConfig{
			ChatMemberConfig: chatMember,
			OnlyIfBanned:     true,
		}
		utils.UnbanChatMember(unbanChatMemberConfig)
	}
}
