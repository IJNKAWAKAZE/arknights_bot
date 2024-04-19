package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"crypto/rand"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"math/big"
)

// QuizHandle 云玩家检测
func QuizHandle(update tgbotapi.Update) error {
	chatId := update.Message.Chat.ID
	userId := update.Message.From.ID
	param := update.Message.CommandArguments()
	key := fmt.Sprintf("quiz:%d", chatId)

	update.Message.Delete()

	if param == "" {
		if utils.RedisIsExists(key) && utils.RedisGet(key) == "stop" {
			sendMessage := tgbotapi.NewMessage(chatId, "云玩家检测功能已关闭！")
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}
	}

	if param == "start" || param == "stop" {
		if bot.Arknights.IsAdmin(chatId, userId) {
			text := ""
			if param == "start" {
				utils.RedisSet(key, "start", 0)
				text = "云玩家检测已开启！"
			} else if param == "stop" {
				utils.RedisSet(key, "stop", 0)
				text = "云玩家检测已关闭！"
			}
			sendMessage := tgbotapi.NewMessage(chatId, text)
			msg, err := bot.Arknights.Send(sendMessage)
			if err != nil {
				return err
			}
			messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
			return nil
		}
		sendMessage := tgbotapi.NewMessage(chatId, "无使用权限！")
		msg, err := bot.Arknights.Send(sendMessage)
		if err != nil {
			return err
		}
		messagecleaner.AddDelQueue(msg.Chat.ID, msg.MessageID, bot.MsgDelDelay)
		return nil
	}

	sendAction := tgbotapi.NewChatAction(chatId, "typing")
	bot.Arknights.Send(sendAction)

	operatorsPool := utils.GetOperators()
	var randNumMap = make(map[int64]struct{})
	var options []utils.Operator
	for i := 0; i < 6; i++ {
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
		skins := operator.Skins
		rsk, _ := rand.Int(rand.Reader, big.NewInt(int64(len(skins))))
		painting := skins[rsk.Int64()].Url
		if param == "h" || param == "ex" {
			painting = skins[0].Url
		}
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

	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: utils.GetImg(correct.ThumbURL)})
	if param == "h" {
		pic := utils.ImgConvert(correct.ThumbURL)
		if pic == nil {
			return nil
		}
		sendPhoto = tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	}
	if param == "ex" {
		pic := utils.CutImg(correct.ThumbURL)
		if pic == nil {
			return nil
		}
		sendPhoto = tgbotapi.NewPhoto(chatId, tgbotapi.FileBytes{Bytes: pic})
	}
	photo, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		log.Printf("发送图片失败：%s，原因：%s", correct.ThumbURL, err.Error())
		return nil
	}
	messagecleaner.AddDelQueue(chatId, photo.MessageID, 300)
	poll := tgbotapi.NewPoll(chatId, "请选择上图干员的正确名字")
	poll.IsAnonymous = false
	poll.Type = "quiz"
	poll.CorrectOptionID = r.Int64()
	var pollOptions []string
	for _, v := range options {
		pollOptions = append(pollOptions, v.Name)
	}
	poll.Options = pollOptions
	p, err := bot.Arknights.Send(poll)
	if err != nil {
		return err
	}
	messagecleaner.AddDelQueue(chatId, p.MessageID, 300)
	return nil
}
