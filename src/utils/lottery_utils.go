package utils

import (
	bot "arknights_bot/config"
	"gorm.io/gorm"
	"time"
)

type GroupLottery struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	GroupName   string    `json:"groupName"`
	GroupNumber int64     `json:"groupNumber"`
	Status      int64     `json:"status"`
	EndTime     time.Time `json:"endTime"`
	CreateTime  time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark      string    `json:"remark"`
}

type GroupLotteryDetail struct {
	Id            string    `json:"id" gorm:"primaryKey"`
	LotteryId     string    `json:"lotteryId"`
	UserName      string    `json:"userName"`
	UserNumber    int64     `json:"userNumber"`
	LotteryNumber int64     `json:"lotteryNumber"`
	Status        int64     `json:"status"`
	CreateTime    time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime    time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark        string    `json:"remark"`
}

// GetAllGroupLottery 查询所有抽奖记录
func GetAllGroupLottery() *gorm.DB {
	return bot.DBEngine.Raw("select * from group_lottery where status in (1, 2)")
}

// GetGroupLottery 查询群组抽奖记录
func GetGroupLottery(chatId int64) *gorm.DB {
	return bot.DBEngine.Raw("select * from group_lottery where group_number = ? and status in (1, 2)", chatId)
}

// GetLotteryDetails 查询抽奖参与列表
func GetLotteryDetails(lotteryId string) *gorm.DB {
	return bot.DBEngine.Raw("select * from group_lottery_detail where lottery_id = ? order by lottery_number", lotteryId)
}

// GetLotteryDetail 查询抽奖详情
func GetLotteryDetail(lotteryId string, lotteryNum int) *gorm.DB {
	return bot.DBEngine.Raw("select * from group_lottery_detail where lottery_id = ? and lottery_number = ?", lotteryId, lotteryNum)
}
