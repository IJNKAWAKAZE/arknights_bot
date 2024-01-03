package account

import "time"

type UserPlayer struct {
	Id         string    `json:"id" gorm:"primaryKey"`
	AccountId  string    `json:"accountId"`
	Uid        string    `json:"uid"`
	ServerName string    `json:"serverName"`
	PlayerName string    `json:"playerName"`
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark     string    `json:"remark"`
}
