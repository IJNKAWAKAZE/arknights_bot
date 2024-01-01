package scheduled

import (
	bot "arknights_bot/bot/init"
	"arknights_bot/bot/modules"
	"arknights_bot/bot/utils"
	"encoding/json"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"time"
)

// DelMsg 删除消息
func DelMsg() func() {
	delMsg := func() {
		var msgObject modules.MsgObject
		msgList := utils.RedisGetList("msgObjects")
		for _, msg := range msgList {
			err := json.Unmarshal([]byte(msg), &msgObject)
			if err != nil {
				log.Println(err)
				return
			}
			t := time.Now()
			if t.Sub(msgObject.CreateTime).Minutes() > msgObject.DelTime {
				delMsg := tgbotapi.NewDeleteMessage(msgObject.ChatId, msgObject.MessageId)
				bot.Arknights.Send(delMsg)
				m, _ := json.Marshal(msgObject)
				utils.RedisDelListItem("msgObjects", string(m))
			}
		}
	}
	return delMsg
}
