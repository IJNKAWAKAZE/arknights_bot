package gatekeeper

import (
	bot "arknights_bot/config"
	"log"
	"time"
)

// verifyStore 待验证入群申请状态的持久化。
// 验证状态原先仅存内存，机器人重启后挂起的入群申请无人处理，用户永久卡死；
// 持久化后可在启动时恢复定时器，过期申请统一拒绝。

const verifyTable = "join_request_verify"

type joinRequestVerify struct {
	ChatID    int64     `gorm:"column:chat_id;primaryKey"`
	UserID    int64     `gorm:"column:user_id;primaryKey"`
	Correct   string    `gorm:"column:correct"`
	MessageID int       `gorm:"column:message_id"`
	ExpireAt  time.Time `gorm:"column:expire_at"`
}

// verifyExpireSlack 持久化条目的过期余量：add 先于发送图片，实际超时为 60 秒，
// 此处留出发送与网络余量
const verifyExpireSlack = 90 * time.Second

func ensureVerifyTable() {
	err := bot.DBEngine.Exec("CREATE TABLE IF NOT EXISTS `join_request_verify` (`chat_id` bigint NOT NULL, `user_id` bigint NOT NULL, `correct` varchar(64) NOT NULL DEFAULT '', `message_id` int NOT NULL DEFAULT 0, `expire_at` datetime NOT NULL, PRIMARY KEY (`chat_id`,`user_id`)) ENGINE = InnoDB CHARACTER SET = utf8mb4 COMMENT = '入群验证待处理状态'").Error
	if err != nil {
		log.Printf("入群验证：创建持久化表失败：%v", err)
	}
}

// verifySave 持久化待验证状态（存在则刷新）
func verifySave(chatId int64, userId int64, correct string) {
	err := bot.DBEngine.Exec("INSERT INTO `join_request_verify` (`chat_id`, `user_id`, `correct`, `message_id`, `expire_at`) VALUES (?, ?, ?, 0, ?) ON DUPLICATE KEY UPDATE `correct` = VALUES(`correct`), `expire_at` = VALUES(`expire_at`)", chatId, userId, correct, time.Now().Add(verifyExpireSlack)).Error
	if err != nil {
		log.Printf("入群验证：持久化验证状态失败：%v", err)
	}
}

// verifySaveMessage 记录验证图片消息 ID，供超时后删除
func verifySaveMessage(chatId int64, userId int64, messageId int) {
	err := bot.DBEngine.Exec("UPDATE `join_request_verify` SET `message_id` = ? WHERE `chat_id` = ? AND `user_id` = ?", messageId, chatId, userId).Error
	if err != nil {
		log.Printf("入群验证：更新验证消息 ID 失败：%v", err)
	}
}

// verifyRemove 清理持久化状态
func verifyRemove(chatId int64, userId int64) {
	err := bot.DBEngine.Exec("DELETE FROM `join_request_verify` WHERE `chat_id` = ? AND `user_id` = ?", chatId, userId).Error
	if err != nil {
		log.Printf("入群验证：清理持久化状态失败：%v", err)
	}
}

// RecoverPendingVerifications 启动时恢复挂起的入群验证：
// 已过期的拒绝入群申请，未过期的重建内存状态并重启 60 秒超时定时器
func RecoverPendingVerifications() {
	ensureVerifyTable()
	var rows []joinRequestVerify
	if err := bot.DBEngine.Table(verifyTable).Find(&rows).Error; err != nil {
		log.Printf("入群验证：读取待恢复验证状态失败：%v", err)
		return
	}
	now := time.Now()
	for _, row := range rows {
		if row.ExpireAt.Before(now) {
			log.Printf("入群验证：恢复时发现过期验证，拒绝用户 %d 加入群 %d", row.UserID, row.ChatID)
			bot.Arknights.DeclineChatJoinRequest(row.ChatID, row.UserID)
			verifyRemove(row.ChatID, row.UserID)
			continue
		}
		verifySet.add(row.UserID, row.ChatID, row.Correct)
		go requestVerify(row.ChatID, row.UserID, row.MessageID)
	}
}
