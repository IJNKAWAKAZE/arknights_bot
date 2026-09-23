package repo

import (
	"arknights_bot/config"
	"fmt"
	"gorm.io/gorm"
	"log"
)

// CheckDB 区分数据库操作失败与查询成功但没有记录的情况。
// 发生错误时，调用方必须立即返回，避免提示成功或修改内存状态。
func CheckDB(result *gorm.DB, chatID int64) error {
	if result.Error == nil {
		return nil
	}
	if _, err := config.Arknights.SendText(chatID, "数据库操作失败，请稍后重试。"); err != nil {
		log.Println("发送数据库错误提示失败:", err)
	}
	return fmt.Errorf("数据库操作失败: %w", result.Error)
}
