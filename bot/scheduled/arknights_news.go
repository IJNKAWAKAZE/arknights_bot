package scheduled

import (
	bot "arknights_bot/bot/init"
	"arknights_bot/bot/utils"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

func BilibiliNews() func() {
	return func() {
		text, pics := utils.ParseBilibiliDynamic()
		if len(text) == 0 {
			return
		}
		groups := utils.GetJoinedGroups()
		if pics == nil {
			for _, group := range groups {
				groupNumber, _ := strconv.ParseInt(group, 10, 64)
				sendMessage := tgbotapi.NewMessage(groupNumber, text)
				bot.Arknights.Send(sendMessage)
			}
			return
		}

		if len(pics) == 1 {
			for _, group := range groups {
				groupNumber, _ := strconv.ParseInt(group, 10, 64)
				sendPhoto := tgbotapi.NewPhoto(groupNumber, tgbotapi.FileURL(pics[0]))
				sendPhoto.Caption = text
				bot.Arknights.Send(sendPhoto)
			}
			return
		}

		for _, group := range groups {
			groupNumber, _ := strconv.ParseInt(group, 10, 64)
			var mediaGroup tgbotapi.MediaGroupConfig
			var media []interface{}
			mediaGroup.ChatID = groupNumber
			for i, pic := range pics {
				var inputPhoto tgbotapi.InputMediaPhoto
				inputPhoto.Media = tgbotapi.FileURL(pic)
				inputPhoto.Type = "photo"
				if i == 0 {
					inputPhoto.Caption = text
				}
				media = append(media, inputPhoto)
			}
			mediaGroup.Media = media
			bot.Arknights.SendMediaGroup(mediaGroup)
		}
	}
}
