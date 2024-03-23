package player

import (
	bot "arknights_bot/config"
	"arknights_bot/plugins/commandoperation"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
		answer := tgbotapi.NewCallbackWithAlert(callbackQuery.ID, "这不是你的角色！")
		bot.Arknights.Send(answer)
		return nil
	} else {
		return callback.Function(playerId)
	}
}
