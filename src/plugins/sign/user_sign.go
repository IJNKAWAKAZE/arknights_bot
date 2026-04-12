package sign

import "time"

type UserSign struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	UserName    string    `json:"userName"`
	UserNumber  int64     `json:"userNumber"`
	NotifyMode  int       `json:"notifyMode"`  // 签到通知模式 0-全部通知 1-仅失败通知 2-仅成功通知
	ApRemind    int       `json:"apRemind"`     // 理智提醒开关 0-关闭 1-开启
	ApThreshold int       `json:"apThreshold"`  // 理智提醒阈值百分比 默认80
	ApNotified  int       `json:"apNotified"`   // 理智提醒是否已通知 0-未通知 1-已通知
	CreateTime  time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark      string    `json:"remark"`
}
