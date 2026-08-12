package gatekeeper

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 入群审计日志：记录用户进群/被拒事件及渠道，用于排查是否因 BUG 未经验证被放入
var auditMu sync.Mutex

func auditJoin(chatId, userId int64, name, channel, detail string) {
	line := time.Now().Format("2006-01-02 15:04:05") +
		fmt.Sprintf(" 群:%d 用户:%d 昵称:%s 渠道:%s 详情:%s", chatId, userId, name, channel, detail)
	// stdout，方便 docker logs / 控制台查看
	log.Println("入群审计：" + line)
	// 落盘 logs/join.log，方便本地审计
	auditMu.Lock()
	defer auditMu.Unlock()
	if err := os.MkdirAll("logs", 0o755); err != nil {
		log.Println("创建日志目录失败:", err)
		return
	}
	f, err := os.OpenFile(filepath.Join("logs", "join.log"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		log.Println("写入入群审计日志失败:", err)
		return
	}
	defer f.Close()
	if _, err := f.WriteString(line + "\n"); err != nil {
		log.Println("写入入群审计日志失败:", err)
	}
}
