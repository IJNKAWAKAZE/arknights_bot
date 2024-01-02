package utils

import "time"

type GroupInvite struct {
	Id           string    `json:"id" gorm:"primaryKey"`
	GroupName    string    `json:"groupName"`
	GroupNumber  string    `json:"groupNumber"`
	UserName     string    `json:"userName"`
	UserNumber   string    `json:"userNumber"`
	MemberName   string    `json:"memberName"`
	MemberNumber string    `json:"memberNumber"`
	CreateTime   time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark       string    `json:"remark"`
	Deleted      int64     `json:"deleted"`
}
