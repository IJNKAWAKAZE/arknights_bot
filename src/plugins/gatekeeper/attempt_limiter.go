package gatekeeper

import (
	"fmt"
	"sync"
	"time"
)

const (
	// maxFailAttempts 连续答错次数上限，超过后限制再次申请入群
	maxFailAttempts = 3
	// blockDuration 超过答错上限后的封禁时长
	blockDuration = 24 * time.Hour
)

// attemptLimiter 记录入群申请答题失败次数，防止通过反复申请穷举答案
type attemptLimiter struct {
	mu      sync.Mutex
	fails   map[string]int
	blocked map[string]time.Time
}

var limiter = attemptLimiter{
	fails:   make(map[string]int),
	blocked: make(map[string]time.Time),
}

// recordFail 记录一次答题失败，达到上限后进入封禁状态，返回累计失败次数
func (l *attemptLimiter) recordFail(chatId int64, userId int64) int {
	l.mu.Lock()
	defer l.mu.Unlock()
	key := limiterKey(chatId, userId)
	l.fails[key]++
	if l.fails[key] >= maxFailAttempts {
		l.blocked[key] = time.Now().Add(blockDuration)
	}
	return l.fails[key]
}

// resetFail 答题通过后清零失败记录与封禁状态
func (l *attemptLimiter) resetFail(chatId int64, userId int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	key := limiterKey(chatId, userId)
	delete(l.fails, key)
	delete(l.blocked, key)
}

// isBlocked 判断用户是否处于封禁状态，封禁过期后自动解除
func (l *attemptLimiter) isBlocked(chatId int64, userId int64) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	key := limiterKey(chatId, userId)
	until, ok := l.blocked[key]
	if !ok {
		return false
	}
	if time.Now().After(until) {
		delete(l.blocked, key)
		delete(l.fails, key)
		return false
	}
	return true
}

func limiterKey(chatId int64, userId int64) string {
	return fmt.Sprintf("%d:%d", chatId, userId)
}
