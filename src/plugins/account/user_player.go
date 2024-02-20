package account

import "time"

type UserPlayer struct {
	Id         string    `json:"id" gorm:"primaryKey"`
	AccountId  string    `json:"accountId"`
	UserName   string    `json:"userName"`
	UserNumber int64     `json:"userNumber"`
	Uid        string    `json:"uid"`
	ServerName string    `json:"serverName"`
	PlayerName string    `json:"playerName"`
	BToken     string    `json:"BToken"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark     string    `json:"remark"`
}
