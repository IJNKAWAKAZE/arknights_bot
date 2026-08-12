package birthday

import (
	"arknights_bot/config"
	"arknights_bot/utils/cache"
	mediautil "arknights_bot/utils/media"
	"arknights_bot/utils/model"
	"arknights_bot/utils/repo"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"time"
)

func Birthday() {
	var operators []model.Operator
	now := time.Now()
	_, month, day := now.Date()
	list := cache.RedisGet(fmt.Sprintf("birthday:%d月%d日", month, day))
	json.Unmarshal([]byte(list), &operators)

	// 存在过生日干员
	template := "很抱歉打断您的浏览🖐\n今天是%s的生日，把这条消息转发在五个群，%s就会从泰拉飞到你家来陪你，我试过了没有用，还会被群友骂舟批，但是今天真的是%s的生日。"
	groups := repo.GetBirthdayGroups()
	for _, group := range groups {
		if month == time.December && day == 24 {
			sendPhoto := tgbotapi.NewPhoto(group, tgbotapi.FilePath("./assets/common/FrostNova.jpg"))
			sendPhoto.Caption = "很抱歉打断您的浏览🖐\n今天是霜星的上岛纪念日，把这条消息转发在五个群，霜星就会上岛来陪你，我试过了没有用，还会被群友骂舟批，但是今天真的是霜星上岛纪念日。"
			config.Arknights.Send(sendPhoto)
		}
		if len(operators) == 1 {
			name := operators[0].Name
			sendPhoto := tgbotapi.NewPhoto(group, tgbotapi.FileBytes{Bytes: mediautil.GetImg(operators[0].Skins[0].Url)})
			sendPhoto.Caption = fmt.Sprintf(template, name, name, name)
			config.Arknights.Send(sendPhoto)
		} else if len(operators) > 1 {
			var mediaGroup tgbotapi.MediaGroupConfig
			var media []interface{}
			mediaGroup.ChatID = group
			var text string
			for i, operator := range operators {
				var inputPhoto tgbotapi.InputMediaPhoto
				inputPhoto.Media = tgbotapi.FileBytes{Bytes: mediautil.GetImg(operator.Skins[0].Url)}
				inputPhoto.Type = "photo"
				text += operator.Name + "、"
				if i == len(operators)-1 {
					text = text[:len(text)-3]
					inputPhoto.Caption = fmt.Sprintf(template, text, text, text)
				}
				media = append(media, inputPhoto)
			}
			mediaGroup.Media = media
			config.Arknights.SendMediaGroup(mediaGroup)
		}
	}
}
