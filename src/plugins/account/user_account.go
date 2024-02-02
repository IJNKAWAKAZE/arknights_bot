package account

import "time"

type UserAccount struct {
	Id              string    `json:"id" gorm:"primaryKey"`
	UserName        string    `json:"userName"`
	UserNumber      int64     `json:"userNumber"`
	HypergryphToken string    `json:"hypergryphToken"`
	BToken          string    `json:"BToken"`
	SklandToken     string    `json:"sklandToken"`
	SklandCred      string    `json:"sklandCred"`
	CreateTime      time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime      time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark          string    `json:"remark"`
}
