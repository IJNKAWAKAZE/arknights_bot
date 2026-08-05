package gatekeeper

import (
	"fmt"
	"sync"
	"time"
)

// approvedMarker 记录刚被机器人批准入群的用户，
// 用于区分「入群申请通过（已答题）」与「链接直入（未答题）」
type approvedMarker struct {
	mu      sync.Mutex
	entries map[string]time.Time
}

var recentlyApproved = approvedMarker{entries: make(map[string]time.Time)}

// approvedMarkerTTL 批准标记的有效时长，超过后视为普通加入
const approvedMarkerTTL = 10 * time.Minute

func (m *approvedMarker) mark(userId int64, chatId int64) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.entries[markerKey(chatId, userId)] = time.Now()
}

func (m *approvedMarker) is(userId int64, chatId int64) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	key := markerKey(chatId, userId)
	t, ok := m.entries[key]
	if !ok {
		return false
	}
	if time.Since(t) > approvedMarkerTTL {
		delete(m.entries, key)
		return false
	}
	return true
}

func markerKey(chatId int64, userId int64) string {
	return fmt.Sprintf("%d:%d", chatId, userId)
}
