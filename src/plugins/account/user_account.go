package account

import "time"

type UserToken struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Content string `json:"content"`
	} `json:"data"`
}

type UserAccount struct {
	Id              string    `json:"id" gorm:"primaryKey"`
	UserName        string    `json:"userName"`
	UserNumber      int64     `json:"userNumber"`
	HypergryphToken string    `json:"hypergryphToken"`
	SklandToken     string    `json:"sklandToken"`
	SklandCred      string    `json:"sklandCred"`
	SklandId        string    `json:"sklandId"`
	ServerName      string    `json:"serverName"`
	CreateTime      time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime      time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark          string    `json:"remark"`
}
