package sign

import "time"

type UserSign struct {
	Id         string    `json:"id" gorm:"primaryKey"`
	UserName   string    `json:"userName"`
	UserNumber int64     `json:"userNumber"`
	NotifyMode int       `json:"notifyMode"` // 签到通知模式 0-全部通知 1-仅失败通知 2-仅成功通知
	CreateTime time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark     string    `json:"remark"`
}
