package apremind

import "time"

type UserApRemind struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	UserName    string    `json:"userName"`
	UserNumber  int64     `json:"userNumber"`
	ApThreshold int       `json:"apThreshold"` // 理智提醒阈值百分比 默认80
	ApNotified  int       `json:"apNotified"`  // 理智提醒是否已通知 0-未通知 1-已通知
	CreateTime  time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark      string    `json:"remark"`
}
