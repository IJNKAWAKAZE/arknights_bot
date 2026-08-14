package messagecleaner

import (
	"arknights_bot/config"
	"arknights_bot/plugins/commandoperation"
	"arknights_bot/utils/cache"
	"encoding/json"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	"log"
	"time"
)

type MsgObject struct {
	ChatId       int64     `json:"chatId"`
	MessageId    int64     `json:"messageId"`
	CreateTime   time.Time `json:"createTime"`
	DelTime      float64   `json:"delTime"`
	FunctionHash string    `json:"functionHash"`
}

// DelMsg 删除消息
func DelMsg() {
	var msgObject MsgObject
	msgList := cache.RedisGetList("msgObjects")
	for _, msg := range msgList {
		err := json.Unmarshal([]byte(msg), &msgObject)
		if err != nil {
			log.Println(err)
			return
		}
		t := time.Now()
		if t.Sub(msgObject.CreateTime).Seconds() > msgObject.DelTime {
			delMsg := tgbotapi.NewDeleteMessage(msgObject.ChatId, msgObject.MessageId)
			config.Arknights.Send(delMsg)
			if len(msgObject.FunctionHash) != 4 { // None 所有的functionHash都是大于4字的
				commandoperation.RemoveCallBack(msgObject.FunctionHash)
			}
			m, _ := json.Marshal(msgObject)
			cache.RedisDelListItem("msgObjects", string(m))
		}
	}
	return
}

// AddDelQueue 添加到删除队列
func AddDelQueue(chatId int64, messageId int64, delTime float64) {
	AddDelQueueFuncHash(chatId, messageId, delTime, "None")
}
func AddDelQueueFuncHash(chatId int64, messageId int64, delTime float64, hash string) {
	var msgObject = MsgObject{
		ChatId:       chatId,
		MessageId:    messageId,
		CreateTime:   time.Now(),
		DelTime:      delTime,
		FunctionHash: hash,
	}
	m, _ := json.Marshal(msgObject)
	cache.RedisSetList("msgObjects", string(m))
}
