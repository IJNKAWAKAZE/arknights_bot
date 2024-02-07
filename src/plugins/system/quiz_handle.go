package system

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/messagecleaner"
	"arknights_bot/utils"
	"crypto/rand"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"math/big"
)

// QuizHandle 云玩家检测
func QuizHandle(update tgbotapi.Update) (bool, error) {
	chatId := update.Message.Chat.ID
	messageId := update.Message.MessageID

	delMsg := tgbotapi.NewDeleteMessage(chatId, messageId)
	bot.Arknights.Send(delMsg)

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
		shipName := operator.Get("name").String()
		skins := operator.Get("skins").Array()
		rsk, _ := rand.Int(rand.Reader, big.NewInt(int64(len(skins))))
		painting := skins[rsk.Int64()].String()
		if painting != "" {
			options = append(options, utils.Operator{
				Name:     shipName,
				ThumbURL: painting,
			})
		} else {
			i--
		}
	}

	r, _ := rand.Int(rand.Reader, big.NewInt(int64(len(options))))
	correct := options[r.Int64()]

	sendPhoto := tgbotapi.NewPhoto(chatId, tgbotapi.FileURL(correct.ThumbURL))
	photo, err := bot.Arknights.Send(sendPhoto)
	if err != nil {
		log.Printf("发送图片失败：%s，原因：%s", correct.ThumbURL, err.Error())
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
	p, _ := bot.Arknights.Send(poll)
	messagecleaner.AddDelQueue(chatId, p.MessageID, 300)
	return true, nil
}
