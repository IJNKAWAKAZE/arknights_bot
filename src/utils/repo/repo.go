package repo

import (
	"arknights_bot/config"
	"arknights_bot/utils/model"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"gorm.io/gorm"
)

// SaveInvite 保存邀请记录
func SaveInvite(message *tgbotapi.Message, member *tgbotapi.User) {
	id, _ := gonanoid.New(32)
	groupInvite := model.GroupInvite{
		Id:           id,
		GroupName:    message.Chat.Title,
		GroupNumber:  message.Chat.ID,
		UserName:     message.From.FullName(),
		UserNumber:   message.From.ID,
		MemberName:   member.FullName(),
		MemberNumber: member.ID,
	}

	config.DBEngine.Table("group_invite").Create(&groupInvite)
}

// SaveJoined 保存入群记录
func SaveJoined(message *tgbotapi.Message) {
	id, _ := gonanoid.New(32)
	groupJoined := model.GroupJoined{
		Id:          id,
		GroupName:   message.Chat.Title,
		GroupNumber: message.Chat.ID,
		News:        0,
		Reg:         -1,
		Welcome:     "",
		Birthday:    0,
		RequestMode: 0,
	}

	config.DBEngine.Table("group_joined").Create(&groupJoined)
}

// GetJoinedByChatId 查询入群记录
func GetJoinedByChatId(chatId int64) *gorm.DB {
	return config.DBEngine.Raw("select * from group_joined where group_number = ? limit 1", chatId)
}

// GetAccountByUserId 查询账号信息
func GetAccountByUserId(userId int64) *gorm.DB {
	return config.DBEngine.Raw("select * from user_account where user_number = ? limit 1", userId)
}

// GetAccountByUserIdAndSklandId 查询账号信息
func GetAccountByUserIdAndSklandId(userId int64, sklandId string) *gorm.DB {
	return config.DBEngine.Raw("select * from user_account where user_number = ? and skland_id = ?", userId, sklandId)
}

// GetAccountByUid 查询账号信息
func GetAccountByUid(userId int64, uid string) *gorm.DB {
	return config.DBEngine.Raw("select t.* from user_account t, user_player t1 where t.id = t1.account_id and t.user_number = ? and t1.uid = ? limit 1", userId, uid)
}

// GetPlayersByUserId 查询绑定角色列表
func GetPlayersByUserId(userId int64) *gorm.DB {
	return config.DBEngine.Raw("select * from user_player where user_number = ?", userId)
}

// GetBPlayersByUserId 查询绑定B服角色列表
func GetBPlayersByUserId(userId int64) *gorm.DB {
	return config.DBEngine.Raw("select * from user_player where user_number = ? and server_name in('b服','bilibili服')", userId)
}

// GetPlayerByUserId 查询绑定角色
func GetPlayerByUserId(userId int64, uid string) *gorm.DB {
	return config.DBEngine.Raw("select * from user_player where user_number = ? and uid = ?", userId, uid)
}

// UpdatePlayerName 更新玩家名称
func UpdatePlayerName(uid, name string) {
	if name == "" {
		return
	}
	config.DBEngine.Exec("update user_player set player_name = ? where uid = ?", name, uid)
}

// GetAutoSign 查询自动签到用户
func GetAutoSign() *gorm.DB {
	return config.DBEngine.Raw("select * from user_sign")
}

// GetAutoSignByUserId 查询自动签到用户
func GetAutoSignByUserId(userId int64) *gorm.DB {
	return config.DBEngine.Raw("select * from user_sign where user_number = ?", userId)
}

// GetApRemindUsers 获取开启理智提醒的用户
func GetApRemindUsers() *gorm.DB {
	return config.DBEngine.Raw("select * from user_ap_remind")
}

// GetApRemindByUserId 查询用户理智提醒设置
func GetApRemindByUserId(userId int64) *gorm.DB {
	return config.DBEngine.Raw("select * from user_ap_remind where user_number = ?", userId)
}

// GetNewsGroups 获取开启消息推送的群组
func GetNewsGroups() []int64 {
	var groups []int64
	config.DBEngine.Raw("select group_number from group_joined where news = 1 group by group_number").Scan(&groups)
	return groups
}

// GetBirthdayGroups 获取开启生日推送的群组
func GetBirthdayGroups() []int64 {
	var groups []int64
	config.DBEngine.Raw("select group_number from group_joined where birthday = 1 group by group_number").Scan(&groups)
	return groups
}

// GetUserGacha 获取角色抽卡记录
func GetUserGacha(userId int64, uid string) *gorm.DB {
	return config.DBEngine.Raw("select * from user_gacha where user_number = ? and uid = ? order by ts desc, pool_order desc", userId, uid)
}

// GetUserPoolCount 获取角色卡池水位
func GetUserPoolCount(userId int64, uid string) *gorm.DB {
	return config.DBEngine.Raw("select pool_name, count(1) pool_count, max(ts) ts from user_gacha where user_number = ? and uid = ? group by pool_name order by ts", userId, uid)
}

// GetAllGroupLottery 查询所有抽奖记录
func GetAllGroupLottery() *gorm.DB {
	return config.DBEngine.Raw("select * from group_lottery where status in (1, 2)")
}

// GetGroupLottery 查询群组抽奖记录
func GetGroupLottery(chatId int64) *gorm.DB {
	return config.DBEngine.Raw("select * from group_lottery where group_number = ? and status in (1, 2)", chatId)
}

// GetLotteryDetails 查询抽奖参与列表
func GetLotteryDetails(lotteryId string) *gorm.DB {
	return config.DBEngine.Raw("select * from group_lottery_detail where lottery_id = ? order by lottery_number", lotteryId)
}

// GetLotteryDetail 查询抽奖详情
func GetLotteryDetail(lotteryId string, lotteryNum int) *gorm.DB {
	return config.DBEngine.Raw("select * from group_lottery_detail where lottery_id = ? and lottery_number = ?", lotteryId, lotteryNum)
}
