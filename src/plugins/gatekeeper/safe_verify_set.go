package gatekeeper

import (
	"fmt"
	"sync"
)

type chatidT = int64
type useridT = int64

// safeCallBack 记录每个群内待验证用户的正确答案
type safeCallBack struct {
	mu  sync.Mutex
	set map[string]string
}

var verifySet = safeCallBack{set: make(map[string]string)}

func (b *safeCallBack) add(userId useridT, chatId chatidT, correct string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.set[verifyKey(chatId, userId)] = correct
}

func (b *safeCallBack) checkExist(userId useridT, chatId chatidT) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	_, ok := b.set[verifyKey(chatId, userId)]
	return ok
}

func (b *safeCallBack) checkExistAndRemove(userId useridT, chatId chatidT) (bool, string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	key := verifyKey(chatId, userId)
	if correct, ok := b.set[key]; ok {
		delete(b.set, key)
		return true, correct
	}
	return false, ""
}

// verifyKey 使用分隔符拼接 (chatId, userId)，避免如 (12,345) 与 (123,45) 产生相同 key
func verifyKey(chatId chatidT, userId useridT) string {
	return fmt.Sprintf("%d:%d", chatId, userId)
}
