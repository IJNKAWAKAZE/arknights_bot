package apremind

import (
	"container/heap"
	"testing"
	"time"
)

// newTestScheduler builds an isolated apScheduler suitable for unit tests.
// No database, no goroutines are started.
func newTestScheduler() *apScheduler {
	s := &apScheduler{
		queue:   make(ApCheckHeap, 0),
		userSet: make(map[int64]*ApCheckItem),
		cache:   make(map[int64]*apUserCache),
		cmdCh:   make(chan apCommand, 10),
	}
	heap.Init(&s.queue)
	return s
}

// ── calcCurrentAp ────────────────────────────────────────────────────────────

func TestCalcCurrentAp_NoElapsed(t *testing.T) {
	// 120 seconds elapsed ÷ 360 s/AP = 0 whole AP gained.
	now := int64(1_000_000)
	got := calcCurrentAp(50, 135, int(now-120), now)
	if got != 50 {
		t.Errorf("want 50, got %d", got)
	}
}

func TestCalcCurrentAp_OnePoint(t *testing.T) {
	// Exactly 360 seconds elapsed → 1 AP gained.
	now := int64(1_000_000)
	got := calcCurrentAp(50, 135, int(now-360), now)
	if got != 51 {
		t.Errorf("want 51, got %d", got)
	}
}

func TestCalcCurrentAp_MultiplePoints(t *testing.T) {
	// 1800 s elapsed ÷ 360 = 5 AP gained.
	now := int64(1_000_000)
	got := calcCurrentAp(50, 135, int(now-1800), now)
	if got != 55 {
		t.Errorf("want 55, got %d", got)
	}
}

func TestCalcCurrentAp_CappedAtMax(t *testing.T) {
	// Would gain 200 AP but must be capped at maxAp=135.
	now := int64(1_000_000)
	got := calcCurrentAp(100, 135, int(now-360*200), now)
	if got != 135 {
		t.Errorf("want 135 (capped), got %d", got)
	}
}

func TestCalcCurrentAp_ZeroTimestamp(t *testing.T) {
	// lastApAddTime == 0 → treat as invalid, no elapsed adjustment.
	now := int64(1_000_000)
	got := calcCurrentAp(50, 135, 0, now)
	if got != 50 {
		t.Errorf("want 50 for zero timestamp, got %d", got)
	}
}

func TestCalcCurrentAp_FutureTimestamp(t *testing.T) {
	// lastApAddTime is in the future → elapsed is negative → no adjustment.
	now := int64(1_000_000)
	future := int(now + 3600)
	got := calcCurrentAp(50, 135, future, now)
	if got != 50 {
		t.Errorf("want 50 for future timestamp, got %d", got)
	}
}

func TestCalcCurrentAp_AlreadyFull(t *testing.T) {
	// current == maxAp already, no time elapsed.
	now := int64(1_000_000)
	got := calcCurrentAp(135, 135, int(now-120), now)
	if got != 135 {
		t.Errorf("want 135, got %d", got)
	}
}

func TestCalcCurrentAp_NegativeCurrent(t *testing.T) {
	// Defensive: if the API ever returns a negative current, result should be 0.
	now := int64(1_000_000)
	got := calcCurrentAp(-5, 135, 0, now)
	if got != 0 {
		t.Errorf("want 0 for negative current, got %d", got)
	}
}

func TestCalcCurrentAp_ZeroMaxAp(t *testing.T) {
	// maxAp == 0 → result should be clamped to 0.
	now := int64(1_000_000)
	got := calcCurrentAp(50, 0, int(now-360), now)
	if got != 0 {
		t.Errorf("want 0 for zero maxAp, got %d", got)
	}
}

func TestCalcCurrentAp_NegativeCurrentWithElapsed(t *testing.T) {
	// Negative current but enough elapsed time to become positive.
	now := int64(1_000_000)
	// -3 + 5 elapsed = 2
	got := calcCurrentAp(-3, 135, int(now-1800), now)
	if got != 2 {
		t.Errorf("want 2 for negative current + elapsed, got %d", got)
	}
}

func TestCalcCurrentAp_NegativeCurrentNotEnoughElapsed(t *testing.T) {
	// Negative current with not enough elapsed time → still negative → clamped to 0.
	now := int64(1_000_000)
	// -10 + 1 elapsed = -9 → 0
	got := calcCurrentAp(-10, 135, int(now-360), now)
	if got != 0 {
		t.Errorf("want 0 for still-negative result, got %d", got)
	}
}

// ── scheduleUser ─────────────────────────────────────────────────────────────

func TestScheduleUser_AddsNewItem(t *testing.T) {
	s := newTestScheduler()
	checkTime := time.Now().Add(5 * time.Minute)

	s.scheduleUser(100, checkTime)

	if s.queue.Len() != 1 {
		t.Fatalf("want 1 item in queue, got %d", s.queue.Len())
	}
	item, ok := s.userSet[100]
	if !ok {
		t.Fatal("user 100 not in userSet")
	}
	if !item.NextCheckTime.Equal(checkTime) {
		t.Errorf("wrong NextCheckTime: %v", item.NextCheckTime)
	}
}

func TestScheduleUser_RescheduleExisting(t *testing.T) {
	s := newTestScheduler()
	now := time.Now()
	s.scheduleUser(100, now.Add(10*time.Minute))
	s.scheduleUser(200, now.Add(5*time.Minute))

	// Reschedule user 100 to become earlier than user 200.
	s.scheduleUser(100, now.Add(1*time.Minute))

	if s.queue.Len() != 2 {
		t.Fatalf("want 2 items, got %d", s.queue.Len())
	}
	// Queue head must now be user 100.
	if s.queue[0].UserNumber != 100 {
		t.Errorf("want user 100 at head after reschedule, got %d", s.queue[0].UserNumber)
	}
}

func TestScheduleUser_MultipleUsers_HeapOrder(t *testing.T) {
	s := newTestScheduler()
	now := time.Now()

	// Insert in non-sorted order.
	delays := []struct {
		user  int64
		delay time.Duration
	}{
		{5, 5 * time.Minute},
		{1, 1 * time.Minute},
		{4, 4 * time.Minute},
		{2, 2 * time.Minute},
		{3, 3 * time.Minute},
	}
	for _, d := range delays {
		s.scheduleUser(d.user, now.Add(d.delay))
	}

	// Pop order must be ascending by user (which correlates with ascending delay).
	for wantUser := int64(1); wantUser <= 5; wantUser++ {
		got := heap.Pop(&s.queue).(*ApCheckItem)
		delete(s.userSet, got.UserNumber)
		if got.UserNumber != wantUser {
			t.Errorf("want user %d, got %d", wantUser, got.UserNumber)
		}
	}
}

// ── removeUser ───────────────────────────────────────────────────────────────

func TestRemoveUser_RemovesItemCorrectly(t *testing.T) {
	s := newTestScheduler()
	s.scheduleUser(100, time.Now().Add(1*time.Minute))
	s.scheduleUser(200, time.Now().Add(2*time.Minute))
	s.scheduleUser(300, time.Now().Add(3*time.Minute))

	s.removeUser(200)

	if s.queue.Len() != 2 {
		t.Fatalf("want 2 items after remove, got %d", s.queue.Len())
	}
	if _, ok := s.userSet[200]; ok {
		t.Error("user 200 should not be in userSet after remove")
	}
	// Remaining users must still be in the set.
	for _, uid := range []int64{100, 300} {
		if _, ok := s.userSet[uid]; !ok {
			t.Errorf("user %d should still be in userSet", uid)
		}
	}
}

func TestRemoveUser_NonExistentIsNoOp(t *testing.T) {
	s := newTestScheduler()
	// Must not panic.
	s.removeUser(999)
	if s.queue.Len() != 0 {
		t.Error("queue should remain empty")
	}
}

func TestRemoveUser_LastItem(t *testing.T) {
	s := newTestScheduler()
	s.scheduleUser(42, time.Now().Add(time.Minute))
	s.removeUser(42)

	if s.queue.Len() != 0 {
		t.Errorf("want empty queue, got %d items", s.queue.Len())
	}
	if _, ok := s.userSet[42]; ok {
		t.Error("user 42 should not be in userSet")
	}
}

// ── dailyCheck ───────────────────────────────────────────────────────────────

func TestDailyCheck_RequeuesUsersNotInQueue(t *testing.T) {
	s := newTestScheduler()

	// User 100: in cache but NOT in queue (was notified, waiting for AP to drop).
	s.cache[100] = &apUserCache{UserNumber: 100, Threshold: 80, ApNotified: 1}
	// User 200: in cache AND in queue (still being monitored).
	s.cache[200] = &apUserCache{UserNumber: 200, Threshold: 80, ApNotified: 0}
	s.scheduleUser(200, time.Now().Add(time.Hour))

	s.dailyCheck()

	if _, ok := s.userSet[100]; !ok {
		t.Error("want user 100 added to queue by dailyCheck")
	}
	if s.queue.Len() != 2 {
		t.Errorf("want 2 items in queue, got %d", s.queue.Len())
	}
}

func TestDailyCheck_DoesNotDuplicateQueuedUsers(t *testing.T) {
	s := newTestScheduler()

	// User 100 is already in both cache and queue.
	s.cache[100] = &apUserCache{UserNumber: 100}
	s.scheduleUser(100, time.Now().Add(time.Hour))

	s.dailyCheck()

	// Must still have exactly 1 item.
	if s.queue.Len() != 1 {
		t.Errorf("want 1 item (no duplicate), got %d", s.queue.Len())
	}
}

func TestDailyCheck_EmptyCacheIsNoOp(t *testing.T) {
	s := newTestScheduler()
	s.dailyCheck()
	if s.queue.Len() != 0 {
		t.Error("queue should remain empty")
	}
}

// ── handleCommand ────────────────────────────────────────────────────────────

func TestHandleCommand_Cancel_RemovesUserAndCache(t *testing.T) {
	s := newTestScheduler()
	s.cache[100] = &apUserCache{UserNumber: 100}
	s.scheduleUser(100, time.Now().Add(time.Hour))

	s.handleCommand(apCommand{typ: cmdCancel, userNumber: 100})

	if _, ok := s.userSet[100]; ok {
		t.Error("user 100 should not be in userSet after cancel")
	}
	if _, ok := s.cache[100]; ok {
		t.Error("user 100 should not be in cache after cancel")
	}
	if s.queue.Len() != 0 {
		t.Errorf("queue should be empty after cancel, got %d", s.queue.Len())
	}
}

func TestHandleCommand_Cancel_NotInQueue_OnlyRemovesCache(t *testing.T) {
	s := newTestScheduler()
	// User is only in cache (already notified, out of queue).
	s.cache[200] = &apUserCache{UserNumber: 200, ApNotified: 1}

	s.handleCommand(apCommand{typ: cmdCancel, userNumber: 200})

	if _, ok := s.cache[200]; ok {
		t.Error("user 200 should not be in cache after cancel")
	}
	if s.queue.Len() != 0 {
		t.Error("queue should remain empty")
	}
}

func TestHandleCommand_DailyCheck_RequeuesNotifiedUsers(t *testing.T) {
	s := newTestScheduler()
	s.cache[100] = &apUserCache{UserNumber: 100, ApNotified: 1}

	s.handleCommand(apCommand{typ: cmdDailyCheck})

	if _, ok := s.userSet[100]; !ok {
		t.Error("want user 100 re-queued by cmdDailyCheck")
	}
}

// ── targetUnix calculation (pure arithmetic) ──────────────────────────────────

// TestTargetUnixCalculation verifies the formula used to predict when AP will
// reach the threshold:
//
//	targetUnix = lastApAddTime + apNeeded * apRecoverySeconds
func TestTargetUnixCalculation(t *testing.T) {
	const apNeeded = 10
	const lastApAddTime = 1_000_000
	expected := int64(lastApAddTime) + int64(apNeeded)*apRecoverySeconds
	got := int64(lastApAddTime) + int64(apNeeded)*int64(apRecoverySeconds)
	if got != expected {
		t.Errorf("targetUnix: want %d, got %d", expected, got)
	}
}

// TestThresholdApCalculation verifies percentage-to-AP conversion.
func TestThresholdApCalculation(t *testing.T) {
	cases := []struct {
		maxAp     int
		threshold int
		want      int
	}{
		{135, 80, 108}, // 135*80/100 = 108
		{135, 50, 67},  // 135*50/100 = 67 (integer division)
		{100, 100, 100},
		{100, 1, 1},
	}
	for _, tc := range cases {
		got := tc.maxAp * tc.threshold / 100
		if got != tc.want {
			t.Errorf("maxAp=%d threshold=%d%%: want %d, got %d",
				tc.maxAp, tc.threshold, tc.want, got)
		}
	}
}
