package apremind

import (
	"container/heap"
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"arknights_bot/config"
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils/repo"
)

const apRecoverySeconds = 360 // 1 AP per 6 minutes
const defaultApThreshold = 80 // default AP threshold percentage

// calcCurrentAp returns the current AP accounting for the time elapsed since the
// last AP recovery tick. It clamps the result to [0, maxAp] and silently skips
// the elapsed calculation when lastApAddTime is invalid (zero or in the future).
func calcCurrentAp(current, maxAp, lastApAddTime int, now int64) int {
	elapsedAp := 0
	if lastApAddTime > 0 {
		elapsed := now - int64(lastApAddTime)
		if elapsed > 0 {
			elapsedAp = int(elapsed / apRecoverySeconds)
		}
	}
	result := current + elapsedAp
	if result < 0 {
		result = 0
	}
	if result > maxAp {
		result = maxAp
	}
	return result
}

// ---------------------------------------------------------------------------
// In-memory cache types (reduce DB queries)
// ---------------------------------------------------------------------------

// apPlayerCache holds cached player credentials.
type apPlayerCache struct {
	Uid             string
	PlayerName      string
	ServerName      string
	HypergryphToken string
	SklandToken     string
	SklandCred      string
}

// apUserCache holds cached user AP reminder settings and player data.
type apUserCache struct {
	UserNumber int64
	Threshold  int
	ApNotified int
	Players    []apPlayerCache
}

// ---------------------------------------------------------------------------
// CSP command types
// ---------------------------------------------------------------------------

type apCommandType int

const (
	cmdSchedule   apCommandType = iota // Schedule/reschedule a user's AP check
	cmdCancel                          // Cancel and remove a user's AP check
	cmdRefresh                         // Refresh a user's cached data
	cmdDailyCheck                      // Re-check all notified users after daily sign
)

type apCommand struct {
	typ        apCommandType
	userNumber int64
}

// ---------------------------------------------------------------------------
// Scheduler (single consumer goroutine owns all mutable state)
// ---------------------------------------------------------------------------

// apScheduler manages the AP check priority queue with CSP architecture.
// All queue and cache mutations happen exclusively in the consumer goroutine.
type apScheduler struct {
	stop     chan struct{}
	done     chan struct{}
	stopOnce sync.Once
	queue    ApCheckHeap
	userSet  map[int64]*ApCheckItem // 记录当前队列中的用户
	cache    map[int64]*apUserCache // 用户数据的内存缓存
	cmdCh    chan apCommand         // 调度指令通道
}

var scheduler *apScheduler

// ---------------------------------------------------------------------------
// Public API (thread-safe – communicate via channel)
// ---------------------------------------------------------------------------

// InitApRemind initializes the AP reminder scheduler and starts the consumer.
// Called once at startup.
func InitApRemind() {
	scheduler = &apScheduler{
		stop: make(chan struct{}), done: make(chan struct{}),
		queue:   make(ApCheckHeap, 0),
		userSet: make(map[int64]*ApCheckItem),
		cache:   make(map[int64]*apUserCache),
		cmdCh:   make(chan apCommand, 100),
	}
	heap.Init(&scheduler.queue)

	// Pre-load all users with AP remind enabled into cache and queue.
	var users []UserApRemind
	res := repo.GetApRemindUsers().Scan(&users)
	if res.RowsAffected > 0 {
		log.Println("初始化理智提醒调度器...")
		for _, user := range users {
			scheduler.loadUserCache(user.UserNumber)
			item := &ApCheckItem{
				UserNumber:    user.UserNumber,
				NextCheckTime: time.Now(),
			}
			heap.Push(&scheduler.queue, item)
			scheduler.userSet[user.UserNumber] = item
		}
		log.Printf("已加载 %d 个理智提醒用户到调度队列", len(users))
	}

	// Start the single consumer goroutine.
	go scheduler.run()
}

// ScheduleNextApCheck adds/reschedules an AP check for a user.
func ScheduleNextApCheck(userNumber int64) {
	if scheduler == nil {
		return
	}
	scheduler.submit(apCommand{typ: cmdSchedule, userNumber: userNumber})
}

// CancelApCheck cancels a user's AP check and removes them from the scheduler.
func CancelApCheck(userNumber int64) {
	if scheduler == nil {
		return
	}
	scheduler.submit(apCommand{typ: cmdCancel, userNumber: userNumber})
}

// RefreshApUserCache refreshes the in-memory cache for a user.
func RefreshApUserCache(userNumber int64) {
	if scheduler == nil {
		return
	}
	scheduler.submit(apCommand{typ: cmdRefresh, userNumber: userNumber})
}

// DailyApCheck re-queues all users not currently in the queue so the consumer
// can detect whether their AP has dropped below threshold. Called after daily sign.
func DailyApCheck() {
	if scheduler == nil {
		return
	}
	scheduler.submit(apCommand{typ: cmdDailyCheck})
}

// ---------------------------------------------------------------------------
// Cache helpers (called from consumer goroutine only – no locks needed)
// ---------------------------------------------------------------------------

func (s *apScheduler) loadUserCache(userNumber int64) {
	var user UserApRemind
	if r := repo.GetApRemindByUserId(userNumber).Scan(&user); r.RowsAffected == 0 {
		delete(s.cache, userNumber)
		return
	}

	threshold := user.ApThreshold
	if threshold == 0 {
		threshold = defaultApThreshold
	}

	var players []account.UserPlayer
	repo.GetPlayersByUserId(userNumber).Scan(&players)

	cached := make([]apPlayerCache, 0, len(players))
	for _, p := range players {
		var ua account.UserAccount
		if r := repo.GetAccountByUid(userNumber, p.Uid).Scan(&ua); r.RowsAffected > 0 {
			cached = append(cached, apPlayerCache{
				Uid:             p.Uid,
				PlayerName:      p.PlayerName,
				ServerName:      ua.ServerName,
				HypergryphToken: ua.HypergryphToken,
				SklandToken:     ua.SklandToken,
				SklandCred:      ua.SklandCred,
			})
		}
	}

	s.cache[userNumber] = &apUserCache{
		UserNumber: userNumber,
		Threshold:  threshold,
		ApNotified: user.ApNotified,
		Players:    cached,
	}
}

// ---------------------------------------------------------------------------
// Consumer loop
// ---------------------------------------------------------------------------

func (s *apScheduler) run() {
	defer close(s.done)
	for {
		var timerCh <-chan time.Time
		var timer *time.Timer

		if s.queue.Len() > 0 {
			delay := time.Until(s.queue[0].NextCheckTime)
			if delay <= 0 {
				s.processNext()
				continue
			}
			timer = time.NewTimer(delay)
			timerCh = timer.C
		}

		// timerCh is nil when queue is empty → that case blocks forever,
		// so we only wake on a command.
		select {
		case <-s.stop:
			if timer != nil {
				timer.Stop()
			}
			return
		case <-timerCh:
			s.processNext()
		case cmd := <-s.cmdCh:
			if timer != nil {
				timer.Stop()
			}
			s.handleCommand(cmd)
		}
	}
}

func (s *apScheduler) handleCommand(cmd apCommand) {
	switch cmd.typ {
	case cmdSchedule:
		s.loadUserCache(cmd.userNumber)
		if _, ok := s.cache[cmd.userNumber]; ok {
			s.scheduleUser(cmd.userNumber, time.Now())
		}
	case cmdCancel:
		s.removeUser(cmd.userNumber)
		delete(s.cache, cmd.userNumber)
	case cmdRefresh:
		s.loadUserCache(cmd.userNumber)
	case cmdDailyCheck:
		s.dailyCheck()
	}
}

// ---------------------------------------------------------------------------
// Queue operations (consumer goroutine only)
// ---------------------------------------------------------------------------

func (s *apScheduler) scheduleUser(userNumber int64, checkTime time.Time) {
	if existing, ok := s.userSet[userNumber]; ok {
		existing.NextCheckTime = checkTime
		heap.Fix(&s.queue, existing.heapIndex)
	} else {
		item := &ApCheckItem{
			UserNumber:    userNumber,
			NextCheckTime: checkTime,
		}
		heap.Push(&s.queue, item)
		s.userSet[userNumber] = item
	}
	log.Printf("理智提醒：用户 %d 下次检查时间：%s",
		userNumber, checkTime.Format("2006-01-02 15:04:05"))
}

func (s *apScheduler) removeUser(userNumber int64) {
	if item, ok := s.userSet[userNumber]; ok {
		heap.Remove(&s.queue, item.heapIndex)
		delete(s.userSet, userNumber)
	}
}

// ---------------------------------------------------------------------------
// Processing (consumer goroutine – serialises all API calls)
// ---------------------------------------------------------------------------

func (s *apScheduler) processNext() {
	if s.queue.Len() == 0 {
		return
	}

	item := heap.Pop(&s.queue).(*ApCheckItem)
	delete(s.userSet, item.UserNumber)

	uc, ok := s.cache[item.UserNumber]
	if !ok {
		s.loadUserCache(item.UserNumber)
		uc, ok = s.cache[item.UserNumber]
		if !ok {
			return
		}
	}

	// Beta-distributed random delay to avoid API rate-limiting.
	delay := betaDelay()
	log.Printf("理智提醒：用户 %d API请求延迟 %.1f秒", item.UserNumber, delay.Seconds())
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-s.stop:
		return
	case <-timer.C:
	}

	s.checkUserAp(uc)
}

func (s *apScheduler) checkUserAp(uc *apUserCache) {
	if len(uc.Players) == 0 {
		return
	}

	threshold := uc.Threshold

	for _, player := range uc.Players {
		select {
		case <-s.stop:
			return
		default:
		}
		var ska skland.Account
		ska.Hypergryph.Token = player.HypergryphToken
		ska.Skland.Token = player.SklandToken
		ska.Skland.Cred = player.SklandCred

		pd, _, err := skland.GetPlayerInfo(player.Uid, ska, player.ServerName)
		if err != nil {
			log.Println("理智提醒：获取角色信息失败:", player.PlayerName, err)
			continue // try next player
		}
		repo.UpdatePlayerName(player.Uid, pd.Status.Name)

		ap := pd.Status.Ap
		maxAp := ap.Max
		if maxAp == 0 {
			continue
		}

		thresholdAp := maxAp * threshold / 100

		// Current AP accounting for time elapsed since last tick.
		currentAp := calcCurrentAp(ap.Current, maxAp, ap.LastApAddTime, time.Now().Unix())
		apPercent := currentAp * 100 / maxAp

		if currentAp >= thresholdAp {
			// ── AP at or above threshold ──
			if uc.ApNotified == 0 {
				_, sendErr := config.Arknights.SendText(uc.UserNumber, fmt.Sprintf(
					"⚡ 理智提醒\n角色 %s 当前理智：%d/%d (%d%%)\n已达到设定阈值 %d%%",
					player.PlayerName, currentAp, maxAp, apPercent, threshold,
				))
				if sendErr != nil {
					log.Println("发送理智通知失败:", sendErr)
					s.scheduleUser(uc.UserNumber, time.Now().Add(time.Minute))
					return
				}
				if err := config.DBEngine.Exec("update user_ap_remind set ap_notified = 1 where user_number = ?", uc.UserNumber).Error; err != nil {
					log.Println("保存理智通知状态失败:", err)
					s.scheduleUser(uc.UserNumber, time.Now().Add(time.Minute))
					return
				}
				uc.ApNotified = 1
				log.Printf("理智提醒：用户 %d 角色 %s 理智已达阈值，已通知", uc.UserNumber, player.PlayerName)
			}
			// Notified – remove from queue; daily check will re-add when AP drops.
			return
		}

		// ── AP below threshold ──
		if uc.ApNotified == 1 {
			if err := config.DBEngine.Exec("update user_ap_remind set ap_notified = 0 where user_number = ?", uc.UserNumber).Error; err != nil {
				log.Println("保存理智通知状态失败:", err)
				s.scheduleUser(uc.UserNumber, time.Now().Add(time.Minute))
				return
			}
			uc.ApNotified = 0
		}

		// Edge case: ap.Current (from API, without elapsed-time adjustment) may already
		// be at/above threshold even though the time-adjusted currentAp fell into the
		// else branch due to rounding or stale LastApAddTime. Fall back to a short retry.
		apNeeded := thresholdAp - ap.Current
		if apNeeded <= 0 {
			s.scheduleUser(uc.UserNumber, time.Now().Add(time.Minute))
			return
		}
		targetUnix := int64(ap.LastApAddTime) + int64(apNeeded)*apRecoverySeconds
		nextTime := time.Unix(targetUnix, 0).Add(30 * time.Second)
		if nextTime.Before(time.Now().Add(time.Minute)) {
			nextTime = time.Now().Add(time.Minute)
		}

		s.scheduleUser(uc.UserNumber, nextTime)
		return // first successful player only
	}

	// All players failed – retry in 30 minutes.
	s.scheduleUser(uc.UserNumber, time.Now().Add(30*time.Minute))
}

// ---------------------------------------------------------------------------
// Daily check – re-queue users who were notified but not in the queue
// ---------------------------------------------------------------------------

func (s *apScheduler) dailyCheck() {
	log.Println("理智提醒：执行每日理智检查...")
	var toCheck []int64
	for uid := range s.cache {
		if _, inQueue := s.userSet[uid]; !inQueue {
			toCheck = append(toCheck, uid)
		}
	}
	for _, uid := range toCheck {
		s.scheduleUser(uid, time.Now())
	}
	log.Printf("理智提醒：每日检查添加了 %d 个用户到队列", len(toCheck))
}

// StopApRemind 取消尚未完成的延迟等待，并等待当前 API 调用结束。
func StopApRemind(ctx context.Context) error {
	if scheduler == nil {
		return nil
	}
	return scheduler.shutdown(ctx)
}
func (s *apScheduler) shutdown(ctx context.Context) error {
	s.stopOnce.Do(func() { close(s.stop) })
	select {
	case <-s.done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (s *apScheduler) submit(cmd apCommand) {
	select {
	case <-s.stop:
		return
	default:
	}
	select {
	case <-s.stop:
	case s.cmdCh <- cmd:
	}
}
