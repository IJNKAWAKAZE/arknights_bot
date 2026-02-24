package gatekeeper

import (
	"fmt"
	"sync"
)

type chatidT = int64
type useridT = int64
type safeCallBack struct {
	mu      sync.Mutex
	set     map[string]bool
	correct string
}

var verifySet = safeCallBack{set: make(map[string]bool)}

func (b *safeCallBack) add(userId useridT, chatId chatidT, correct string) {
	b.mu.Lock()
	val := fmt.Sprintf("%d%d", chatId, userId)
	b.set[val] = true
	b.correct = correct
	b.mu.Unlock()
}

func (b *safeCallBack) checkExist(userId useridT, chatId chatidT) bool {
	defer b.mu.Unlock()
	b.mu.Lock()
	val := fmt.Sprintf("%d%d", chatId, userId)
	if b._checkVal(val) {
		return true
	}
	return false
}

func (b *safeCallBack) checkExistAndRemove(userId useridT, chatId chatidT) (bool, string) {
	defer b.mu.Unlock()
	b.mu.Lock()
	val := fmt.Sprintf("%d%d", chatId, userId)
	if b._checkVal(val) {
		delete(b.set, val)
		return true, b.correct
	}
	return false, ""
}

func (b *safeCallBack) _checkVal(val string) bool {
	_, ok := b.set[val]
	return ok
}
