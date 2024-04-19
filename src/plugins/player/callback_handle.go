package player

import (
	"arknights_bot/plugins/commandoperation"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"strings"
)

// PlayerData 角色数据
func PlayerData(callBack tgbotapi.Update) error {
	callbackQuery := callBack.CallbackQuery
	data := callBack.CallbackData()
	d := strings.Split(data, ",")
	if len(d) != 3 {
		log.Printf("invalid callback ")
	}
	playerId := d[2]
	callback, ok := commandoperation.GetCallback(d[1])
	if !ok {
		if !ok {
			log.Printf("Unable to Call Back %s", d[1])
		}
		return nil
	}
	userId := callbackQuery.From.ID
	clickUserId := callback.UserId
	if userId != clickUserId {
		callbackQuery.Answer(true, "这不是你的角色！")
		return nil
	} else {
		callbackQuery.Answer(false, "正在渲染图片请勿重复点击")
		return callback.Function(playerId)
	}
}
