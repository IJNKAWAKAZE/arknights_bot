package cron

import "time"

type MsgObject struct {
	ChatId     int64     `json:"chatId"`
	MessageId  int       `json:"messageId"`
	CreateTime time.Time `json:"createTime"`
	DelTime    float64   `json:"delTime"`
}
