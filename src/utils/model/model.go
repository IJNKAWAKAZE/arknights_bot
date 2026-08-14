package model

import "time"

type Operator struct {
	Name         string     `json:"name"`         // 名字
	NameEn       string     `json:"nameEn"`       // 英文名
	NameJa       string     `json:"nameJp"`       // 日文名
	Pinyin       [][]string `json:"pinyin"`       // 拼音变体数组
	Code         string     `json:"code"`         // 编号
	Race         string     `json:"race"`         // 种族
	Profession   string     `json:"profession"`   // 职业
	ProfessionZH string     `json:"professionZH"` // 职业
	Rarity       int        `json:"rarity"`       // 稀有度
	Avatar       string     `json:"avatar"`       // 头像
	ThumbURL     string     `json:"thumbURL"`     // 半身像
	Skins        []Skin     `json:"skins"`        // 皮肤
	HP           string     `json:"hp"`           // 生命值
	ATK          string     `json:"atk"`          // 攻击
	DEF          string     `json:"def"`          // 防御
	Res          string     `json:"res"`          // 法抗
	ReDeploy     string     `json:"reDeploy"`     // 再部署时间
	Cost         string     `json:"cost"`         // 费用
	Block        string     `json:"block"`        // 阻挡数
	Interval     string     `json:"interval"`     // 攻击间隔
	Sex          string     `json:"sex"`          // 性别
	Position     string     `json:"position"`     // 部署位
	Logo         string     `json:"logo"`         // 所属
	ObtainMethod string     `json:"obtainMethod"` // 获取方式
	Tags         string     `json:"tags"`         // 标签
}

type Skin struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type GroupInvite struct {
	Id           string    `json:"id" gorm:"primaryKey"`
	GroupName    string    `json:"groupName"`
	GroupNumber  int64     `json:"groupNumber"`
	UserName     string    `json:"userName"`
	UserNumber   int64     `json:"userNumber"`
	MemberName   string    `json:"memberName"`
	MemberNumber int64     `json:"memberNumber"`
	CreateTime   time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark       string    `json:"remark"`
}

type GroupJoined struct {
	Id          string    `json:"id" gorm:"primaryKey"`
	GroupName   string    `json:"groupName"`
	GroupNumber int64     `json:"groupNumber"`
	News        int64     `json:"news"`
	Reg         int64     `json:"reg"`
	Welcome     string    `json:"welcome"`
	Birthday    int64     `json:"birthday"`
	RequestMode int64     `json:"requestMode"`
	CreateTime  time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime  time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark      string    `json:"remark"`
}

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
