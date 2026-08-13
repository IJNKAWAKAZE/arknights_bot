package gatekeeper

import (
	"arknights_bot/config"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// 入群审计日志：记录用户进群/被拒事件及渠道，用于排查是否因 BUG 未经验证被放入
var auditMu sync.Mutex

// chatTitleCache 缓存群名，避免每条审计日志都请求一次 Telegram API
var chatTitleMu sync.Mutex
var chatTitleCache = make(map[int64]string)

// getChatTitle 获取群标题，失败时返回空串（日志中降级只显示群 id）
func getChatTitle(chatId int64) string {
	chatTitleMu.Lock()
	defer chatTitleMu.Unlock()
	if title, ok := chatTitleCache[chatId]; ok {
		return title
	}
	chat, err := config.Arknights.GetChatInfo(chatId)
	if err != nil {
		return ""
	}
	chatTitleCache[chatId] = chat.Title
	return chat.Title
}

func auditJoin(chatId, userId int64, name, channel, detail string) {
	chatPart := fmt.Sprintf("群:%d", chatId)
	if title := getChatTitle(chatId); title != "" {
		chatPart = fmt.Sprintf("群:%s(%d)", title, chatId)
	}
	line := time.Now().Format("2006-01-02 15:04:05") +
		fmt.Sprintf(" %s 用户:%d 昵称:%s 渠道:%s 详情:%s", chatPart, userId, name, channel, detail)
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
