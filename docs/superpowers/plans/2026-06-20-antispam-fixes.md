# Antispam Guest Spam Fixes Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 修复 guest spam 的活跃用户时间窗口、restore 恢复逻辑和投票通过后的删消息错误处理，并补齐单测与集成测试覆盖。

**Architecture:** 保留现有 antispam 的 Redis 缓存 + MySQL 落盘结构，只把活跃用户从共享 TTL 的 `set` 调整为按时间戳维护的 Redis `sorted set`。Telegram 动作层通过扩展 `telegramGateway` 增加 `UnbanChatMember`，并把 guest spam 消息删除统一收敛到一个带日志的 helper，保证恢复逻辑和回调文案与真实状态一致。

**Tech Stack:** Go, GORM, Redis v8, telegram-bot-api v1.0.12, GitHub Actions

---

## File Structure

**Files:**
- Modify: `src/plugins/antispam/keys.go`
  作用：定义 Redis key 和活跃窗口辅助 key/函数入口。
- Modify: `src/plugins/antispam/cache.go`
  作用：维护活跃用户缓存、冷启动回填、投票活跃资格和计数逻辑。
- Modify: `src/plugins/antispam/telegram_gateway.go`
  作用：为 antispam 提供 Telegram 动作抽象，补充恢复封禁所需接口。
- Modify: `src/plugins/antispam/telegram_gateway_test.go`
  作用：扩展 fake gateway，记录 unban 调用和错误注入。
- Modify: `src/plugins/antispam/message.go`
  作用：统一 guest spam 删消息 helper，并让自动低信誉路径共享它。
- Modify: `src/plugins/antispam/commands.go`
  作用：修复 restore 路径和投票通过路径的动作及文案。
- Modify: `src/plugins/antispam/guest_bot_spam_test.go`
  作用：新增/调整 unit tests，覆盖新 helper、restore 顺序和回调文案。
- Modify: `src/plugins/antispam/guest_bot_spam_integration_test.go`
  作用：验证 sorted set 活跃窗口、冷启动恢复、restore 对 ban 生效、vote delete failure 文案和日志。

---

### Task 1: Convert Active Users To Redis Sorted Set

**Files:**
- Modify: `src/plugins/antispam/keys.go`
- Modify: `src/plugins/antispam/cache.go`
- Test: `src/plugins/antispam/guest_bot_spam_test.go`
- Test: `src/plugins/antispam/guest_bot_spam_integration_test.go`

- [ ] **Step 1: Write the failing unit tests for active-user expiry semantics**

在 `src/plugins/antispam/guest_bot_spam_test.go` 新增以下测试，先定义我们想要的 sorted set 行为：

```go
func TestActiveUsersExpireIndividually(t *testing.T) {
	setupGuestSpamIntegrationRedisOnly(t)

	now := time.Now()
	writeActiveUserScore(-100100, 1001, now.Add(-2*time.Minute))
	writeActiveUserScore(-100100, 1002, now.Add(-12*time.Minute))

	if got := ActiveUserCount(-100100); got != 1 {
		t.Fatalf("active users = %d, want 1 after pruning expired members", got)
	}
	if !IsActiveUser(-100100, 1001) {
		t.Fatal("fresh user should stay active")
	}
	if IsActiveUser(-100100, 1002) {
		t.Fatal("expired user should be inactive")
	}
}

func TestTrackActivityHandleRefreshesSortedSetScore(t *testing.T) {
	setupGuestSpamIntegrationRedisOnly(t)

	if err := TrackActivityHandle(tgbotapi.Update{Message: trackableMessage(880001, "hello")}); err != nil {
		t.Fatalf("track activity: %v", err)
	}
	if got := ActiveUserCount(integrationChatID); got != 1 {
		t.Fatalf("active users = %d, want 1", got)
	}
	if !IsActiveUser(integrationChatID, 880001) {
		t.Fatal("tracked user should be active")
	}
}
```

同时在 `src/plugins/antispam/guest_bot_spam_integration_test.go` 增加一个失败用例，验证冷启动不会把过期成员重新放回活跃集合：

```go
func TestGuestSpamIntegrationLoadCacheOnlyRestoresFreshActiveUsers(t *testing.T) {
	db := setupGuestSpamIntegration(t)

	fresh := time.Now().Add(-2 * time.Minute)
	stale := time.Now().Add(-20 * time.Minute)
	if err := db.Create(&MemberRisk{
		ID:            riskID(integrationChatID, 9001),
		ChatID:        integrationChatID,
		UserID:        9001,
		UserName:      "Fresh",
		FirstSeenAt:   startOfDay(time.Now().AddDate(0, 0, -5)),
		LastMessageAt: fresh,
	}).Error; err != nil {
		t.Fatalf("create fresh risk: %v", err)
	}
	if err := db.Create(&MemberRisk{
		ID:            riskID(integrationChatID, 9002),
		ChatID:        integrationChatID,
		UserID:        9002,
		UserName:      "Stale",
		FirstSeenAt:   startOfDay(time.Now().AddDate(0, 0, -5)),
		LastMessageAt: stale,
	}).Error; err != nil {
		t.Fatalf("create stale risk: %v", err)
	}

	clearGuestSpamRedis(t)
	if err := LoadCacheFromDB(); err != nil {
		t.Fatalf("load cache from db: %v", err)
	}
	if got := ActiveUserCount(integrationChatID); got != 1 {
		t.Fatalf("active users after reload = %d, want 1", got)
	}
	if !IsActiveUser(integrationChatID, 9001) || IsActiveUser(integrationChatID, 9002) {
		t.Fatal("reload should keep fresh user only")
	}
}
```

- [ ] **Step 2: Run the focused tests to verify they fail against the current set-based implementation**

Run:

```powershell
go test ./plugins/antispam -run "TestActiveUsersExpireIndividually|TestTrackActivityHandleRefreshesSortedSetScore" -count=1
go test -tags=integration ./plugins/antispam -run TestGuestSpamIntegrationLoadCacheOnlyRestoresFreshActiveUsers -count=1
```

Expected:

- 单测至少一个失败，因为当前实现不会按用户独立过期
- 集成测试失败，因为 `LoadCacheFromDB` 会按旧 set 逻辑恢复所有窗口内外混合状态

- [ ] **Step 3: Implement sorted-set active user tracking in keys and cache**

在 `src/plugins/antispam/keys.go` 和 `src/plugins/antispam/cache.go` 做最小实现，核心代码形态如下：

```go
func activeUsersKey(chatID int64) string {
	return fmt.Sprintf("%s:active:%d", redisPrefix, chatID)
}

func activeWindowCutoff(now time.Time) float64 {
	return float64(now.Add(-activeWindowTTL).Unix())
}

func pruneActiveUsers(chatID int64, now time.Time) {
	if bot.GoRedis == nil {
		return
	}
	key := activeUsersKey(chatID)
	if err := bot.GoRedis.ZRemRangeByScore(redisCtx, key, "-inf", fmt.Sprintf("%f", activeWindowCutoff(now))).Err(); err != nil {
		log.Printf("guest spam: prune active users failed: %v", err)
	}
}

func RecordMessageActivity(chatID int64, userID int64, userName string) {
	// 保留现有 risk/activity 逻辑
	// ...
	now := time.Now()
	key := activeUsersKey(chatID)
	pruneActiveUsers(chatID, now)
	if err := bot.GoRedis.ZAdd(redisCtx, key, &redis.Z{
		Score:  float64(now.Unix()),
		Member: strconv.FormatInt(userID, 10),
	}).Err(); err != nil {
		log.Printf("guest spam: update active user failed: %v", err)
	}
	bot.GoRedis.Expire(redisCtx, key, activeWindowTTL)
}

func ActiveUserCount(chatID int64) int {
	if bot.GoRedis == nil {
		return 0
	}
	pruneActiveUsers(chatID, time.Now())
	count, err := bot.GoRedis.ZCard(redisCtx, activeUsersKey(chatID)).Result()
	if err != nil {
		log.Printf("guest spam: count active users failed: %v", err)
		return 0
	}
	return int(count)
}

func IsActiveUser(chatID, userID int64) bool {
	if bot.GoRedis == nil || userID == 0 {
		return false
	}
	pruneActiveUsers(chatID, time.Now())
	score, err := bot.GoRedis.ZScore(redisCtx, activeUsersKey(chatID), strconv.FormatInt(userID, 10)).Result()
	if err == redis.Nil {
		return false
	}
	if err != nil {
		log.Printf("guest spam: read active user failed: %v", err)
		return false
	}
	return score >= activeWindowCutoff(time.Now())
}
```

同时修改 `LoadCacheFromDB()` 的活跃回填逻辑：

```go
for _, risk := range risks {
	setJSON(memberRiskKey(risk.ChatID, risk.UserID), risk, riskTTL)
	if !risk.LastMessageAt.IsZero() && time.Since(risk.LastMessageAt) <= activeWindowTTL {
		if err := bot.GoRedis.ZAdd(redisCtx, activeUsersKey(risk.ChatID), &redis.Z{
			Score:  float64(risk.LastMessageAt.Unix()),
			Member: strconv.FormatInt(risk.UserID, 10),
		}).Err(); err != nil {
			log.Printf("guest spam: restore active user failed: %v", err)
		}
		bot.GoRedis.Expire(redisCtx, activeUsersKey(risk.ChatID), activeWindowTTL)
	}
}
```

为测试提供一个只初始化 Redis 的辅助函数，避免 unit test 因 MySQL 耦合过深：

```go
func setupGuestSpamIntegrationRedisOnly(t *testing.T) {
	t.Helper()
	setupGuestSpamIntegration(t)
	clearGuestSpamTables(t, bot.DBEngine)
}
```

- [ ] **Step 4: Run the focused tests to verify the sorted-set implementation passes**

Run:

```powershell
go test ./plugins/antispam -run "TestActiveUsersExpireIndividually|TestTrackActivityHandleRefreshesSortedSetScore" -count=1
go test -tags=integration ./plugins/antispam -run "TestGuestSpamIntegrationLoadCacheOnlyRestoresFreshActiveUsers|TestGuestSpamIntegrationTrackActivityHandle|TestGuestSpamIntegrationActivityDateSyncAndReload" -count=1
```

Expected:

- 所有上述测试 PASS
- 不再依赖共享 TTL 的 `SAdd` / `SCard` / `SIsMember` 活跃模型

- [ ] **Step 5: Commit the sorted-set active user fix**

```bash
git add src/plugins/antispam/keys.go src/plugins/antispam/cache.go src/plugins/antispam/guest_bot_spam_test.go src/plugins/antispam/guest_bot_spam_integration_test.go
git commit -m "fix: track guest spam active users with sorted sets"
```

---

### Task 2: Add Full Restore Flow With Unban + Unrestrict

**Files:**
- Modify: `src/plugins/antispam/telegram_gateway.go`
- Modify: `src/plugins/antispam/telegram_gateway_test.go`
- Modify: `src/plugins/antispam/commands.go`
- Test: `src/plugins/antispam/guest_bot_spam_test.go`
- Test: `src/plugins/antispam/guest_bot_spam_integration_test.go`

- [ ] **Step 1: Write the failing restore tests**

在 `src/plugins/antispam/guest_bot_spam_test.go` 增加 unit tests，先把调用顺序和失败语义钉住：

```go
func TestRestoreCallerUnbansBeforeClearingRestrictions(t *testing.T) {
	fake := useFakeTelegram(t)
	err := restoreCaller(integrationChatID, integrationCallerID, commandMessage("/guest_spam_log restore 880001", 8001))
	if err != nil {
		t.Fatalf("restore caller: %v", err)
	}
	if len(fake.unbans) != 1 {
		t.Fatalf("unban calls = %d, want 1", len(fake.unbans))
	}
	if len(fake.restricts) != 1 {
		t.Fatalf("restrict calls = %d, want 1", len(fake.restricts))
	}
}

func TestRestoreCallerStopsWhenUnbanFails(t *testing.T) {
	fake := useFakeTelegram(t)
	fake.unbanErr = errTelegram()
	err := restoreCaller(integrationChatID, integrationCallerID, commandMessage("/guest_spam_log restore 880001", 8001))
	if err == nil {
		t.Fatal("restore should fail when unban fails")
	}
	if len(fake.restricts) != 0 {
		t.Fatalf("restrict calls = %d, want 0 after unban failure", len(fake.restricts))
	}
}
```

在 `src/plugins/antispam/guest_bot_spam_integration_test.go` 为已 ban 用户增加恢复验证：

```go
func TestGuestSpamIntegrationRestoreCallerUnbansAndClearsWarnings(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)

	AddWarning(integrationChatID, integrationCallerID, "Caller")
	if err := GuestSpamLogHandle(tgbotapi.Update{Message: commandMessage("/guest_spam_log restore 880001", 8001)}); err != nil {
		t.Fatalf("restore command: %v", err)
	}
	if len(fake.unbans) != 1 {
		t.Fatalf("unban calls = %d, want 1", len(fake.unbans))
	}
	if len(fake.restricts) != 1 {
		t.Fatalf("restrict calls = %d, want 1", len(fake.restricts))
	}
}
```

- [ ] **Step 2: Run the restore-focused tests to verify they fail**

Run:

```powershell
go test ./plugins/antispam -run "TestRestoreCallerUnbansBeforeClearingRestrictions|TestRestoreCallerStopsWhenUnbanFails" -count=1
go test -tags=integration ./plugins/antispam -run TestGuestSpamIntegrationRestoreCallerUnbansAndClearsWarnings -count=1
```

Expected:

- 单测失败，因为当前 fake gateway 和真实 gateway 都没有 `UnbanChatMember`
- 集成测试失败，因为 restore 仍只有 unrestrict 逻辑

- [ ] **Step 3: Extend telegram gateway and update restoreCaller**

修改 `src/plugins/antispam/telegram_gateway.go`，扩展接口和 live 实现：

```go
type telegramGateway interface {
	Send(tgbotapi.Chattable) (tgbotapi.Message, error)
	Request(tgbotapi.Chattable) (*tgbotapi.APIResponse, error)
	DeleteMessage(chatID int64, messageID int) (*tgbotapi.APIResponse, error)
	AnswerCallback(callbackID string, showAlert bool, text string) (*tgbotapi.APIResponse, error)
	DeleteCallbackMessage(callback *tgbotapi.CallbackQuery) (*tgbotapi.APIResponse, error)
	BanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error)
	UnbanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error)
	RestrictChatMember(chatID, userID int64, permissions string) (*tgbotapi.APIResponse, error)
	IsAdmin(chatID, userID int64) bool
	QueueDelete(chatID int64, messageID int, delay float64)
}

func (liveTelegramGateway) UnbanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error) {
	return bot.Arknights.UnbanChatMember(chatID, userID)
}
```

修改 `src/plugins/antispam/telegram_gateway_test.go`，给 fake 增加记录和错误注入：

```go
type fakeTelegramGateway struct {
	// ...
	unbans        []banCall
	unbanErr      error
}

func (fake *fakeTelegramGateway) UnbanChatMember(chatID, userID int64) (*tgbotapi.APIResponse, error) {
	fake.unbans = append(fake.unbans, banCall{chatID: chatID, userID: userID})
	if fake.unbanErr != nil {
		return nil, fake.unbanErr
	}
	return &tgbotapi.APIResponse{Ok: true}, nil
}
```

修改 `src/plugins/antispam/commands.go` 的 `restoreCaller`：

```go
func restoreCaller(chatID, userID int64, message *tgbotapi.Message) error {
	if _, err := guestSpamTelegram.UnbanChatMember(chatID, userID); err != nil {
		AddLog(SpamLog{
			ChatID:       chatID,
			ChatName:     message.Chat.Title,
			CallerUserID: userID,
			Action:       ActionRestoreCaller,
			Reason:       ReasonAdminRestore,
			Detail:       "unban failed: " + err.Error(),
		})
		return err
	}
	if _, err := guestSpamTelegram.RestrictChatMember(chatID, userID, tgbotapi.AllPermissions); err != nil {
		AddLog(SpamLog{
			ChatID:       chatID,
			ChatName:     message.Chat.Title,
			CallerUserID: userID,
			Action:       ActionRestoreCaller,
			Reason:       ReasonAdminRestore,
			Detail:       "unrestrict failed: " + err.Error(),
		})
		return err
	}
	RestoreCallerState(chatID, userID, message)
	return sendTempMessage(chatID, message.MessageID, "已恢复该用户并清除 guest spam 警告。")
}
```

- [ ] **Step 4: Run the restore-focused tests to verify the new flow passes**

Run:

```powershell
go test ./plugins/antispam -run "TestRestoreCallerUnbansBeforeClearingRestrictions|TestRestoreCallerStopsWhenUnbanFails|TestCommandAndCallbackGuards" -count=1
go test -tags=integration ./plugins/antispam -run "TestGuestSpamIntegrationRestoreCallerUnbansAndClearsWarnings|TestGuestSpamIntegrationGuestSpamLogHandlePaths|TestGuestSpamIntegrationRestoreClearsWarnings" -count=1
```

Expected:

- 所有 restore 相关测试 PASS
- 失败路径不再清除 warning 状态
- fake gateway 能记录 unban 和 unrestrict 两段动作

- [ ] **Step 5: Commit the restore fix**

```bash
git add src/plugins/antispam/telegram_gateway.go src/plugins/antispam/telegram_gateway_test.go src/plugins/antispam/commands.go src/plugins/antispam/guest_bot_spam_test.go src/plugins/antispam/guest_bot_spam_integration_test.go
git commit -m "fix: fully restore guest spam callers"
```

---

### Task 3: Unify Guest Spam Message Deletion And Surface Failures

**Files:**
- Modify: `src/plugins/antispam/message.go`
- Modify: `src/plugins/antispam/commands.go`
- Test: `src/plugins/antispam/guest_bot_spam_test.go`
- Test: `src/plugins/antispam/guest_bot_spam_integration_test.go`

- [ ] **Step 1: Write the failing tests for delete helper and vote-pass callback copy**

在 `src/plugins/antispam/guest_bot_spam_test.go` 中加入以下测试：

```go
func TestDeleteGuestMessageLogsFailure(t *testing.T) {
	setupGuestSpamIntegrationRedisOnly(t)
	fake := useFakeTelegram(t)
	fake.deleteErr = errTelegram()

	err := deleteGuestMessageWithLog(guestMessage(2001, 1001), ReasonLowTrust)
	if err == nil {
		t.Fatal("delete helper should return error")
	}
	if !hasLogAction(RecentLogs(-100100, 10), ActionDeleteFailed) {
		t.Fatalf("logs = %+v, want delete_failed", RecentLogs(-100100, 10))
	}
}

func TestApplyVotePassedShowsDeleteFailureMessage(t *testing.T) {
	setupGuestSpamIntegrationRedisOnly(t)
	fake := useFakeTelegram(t)
	fake.deleteErr = errTelegram()

	vote := SpamVote{
		ID:                "vote-failure",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         601,
		GuestBotID:        993201,
		GuestBotName:      "Pass Bot",
		GuestBotUserName:  "pass_bot",
		RequiredVoteCount: 1,
		VoteScore:         1,
	}
	err := applyVotePassed(vote, voteCallback("pass-cb", "guestspam_vote,vote,vote-failure", 7002))
	if err == nil {
		t.Fatal("vote pass should surface delete failure")
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "投票通过，已拉黑，但删除消息失败，请管理员检查权限" {
		t.Fatalf("callbacks = %+v", fake.callbacks)
	}
}
```

并在 `src/plugins/antispam/guest_bot_spam_integration_test.go` 中补投票通过但删消息失败的集成测试：

```go
func TestGuestSpamIntegrationVotePassDeleteFailureIsLoggedAndReported(t *testing.T) {
	setupGuestSpamIntegration(t)
	fake := useFakeTelegram(t)
	fake.deleteErr = errTelegram()
	bot.GoRedis.ZAdd(redisCtx, activeUsersKey(integrationChatID), &redis.Z{Score: float64(time.Now().Unix()), Member: "7001"})
	bot.GoRedis.ZAdd(redisCtx, activeUsersKey(integrationChatID), &redis.Z{Score: float64(time.Now().Unix()), Member: "7002"})

	SaveVote(SpamVote{
		ID:                "pass-vote-failure",
		ChatID:            integrationChatID,
		ChatName:          "Guest Spam Test",
		MessageID:         611,
		GuestBotID:        993211,
		GuestBotName:      "Pass Bot Failure",
		GuestBotUserName:  "pass_bot_failure",
		ExpiresAt:         time.Now().Add(time.Hour),
		RequiredVoteCount: 2,
		Voters:            []int64{7001},
	})
	if err := SpamVoteCallback(tgbotapi.Update{CallbackQuery: voteCallback("pass-failure-cb", "guestspam_vote,vote,pass-vote-failure", 7002)}); err != nil {
		t.Fatalf("pass failure callback: %v", err)
	}
	if !hasLogAction(RecentLogs(integrationChatID, 10), ActionDeleteFailed) {
		t.Fatalf("logs = %+v, want delete_failed", RecentLogs(integrationChatID, 10))
	}
	if len(fake.callbacks) != 1 || fake.callbacks[0].text != "投票通过，已拉黑，但删除消息失败，请管理员检查权限" {
		t.Fatalf("callbacks = %+v", fake.callbacks)
	}
}
```

- [ ] **Step 2: Run the deletion-focused tests to verify they fail**

Run:

```powershell
go test ./plugins/antispam -run "TestDeleteGuestMessageLogsFailure|TestApplyVotePassedShowsDeleteFailureMessage" -count=1
go test -tags=integration ./plugins/antispam -run TestGuestSpamIntegrationVotePassDeleteFailureIsLoggedAndReported -count=1
```

Expected:

- 单测失败，因为 `deleteGuestMessageWithLog` 当前不返回错误，`applyVotePassed` 也没有失败文案分支
- 集成测试失败，因为回调文案仍然固定为“已拉黑并删除消息”

- [ ] **Step 3: Refactor deletion into a shared helper that returns errors**

在 `src/plugins/antispam/message.go` 中把 helper 改成返回错误：

```go
func deleteGuestMessageWithLog(message *tgbotapi.Message, reason string) error {
	if message == nil || message.Chat == nil {
		return nil
	}
	if _, err := guestSpamTelegram.DeleteMessage(message.Chat.ID, message.MessageID); err != nil {
		AddLog(logFromMessage(message, ActionDeleteFailed, reason, err.Error()))
		log.Printf("guest spam: delete message failed: %v", err)
		return err
	}
	AddLog(logFromMessage(message, ActionDeleteMessage, reason, "deleted guest bot message"))
	return nil
}

func deleteGuestMessageByFields(chatID int64, messageID int, logItem SpamLog) error {
	if _, err := guestSpamTelegram.DeleteMessage(chatID, messageID); err != nil {
		logItem.Action = ActionDeleteFailed
		logItem.Detail = err.Error()
		AddLog(logItem)
		return err
	}
	logItem.Action = ActionDeleteMessage
	logItem.Detail = "deleted guest bot message"
	AddLog(logItem)
	return nil
}
```

在 `GuestBotSpamHandle` 里保留当前“不向上冒泡 Telegram delete error”的策略，但改成显式忽略返回值：

```go
if decision.DeleteMessage {
	_ = deleteGuestMessageWithLog(message, decision.Reason)
}
```

在 `src/plugins/antispam/commands.go` 中把 `applyVotePassed` 改成返回 `error` 并按删除结果设置 callback：

```go
func applyVotePassed(vote SpamVote, callback *tgbotapi.CallbackQuery) error {
	ApplyVotePassedState(vote)
	err := deleteGuestMessageByFields(vote.ChatID, vote.MessageID, SpamLog{
		ChatID:       vote.ChatID,
		ChatName:     vote.ChatName,
		MessageID:    vote.MessageID,
		GuestBotID:   vote.GuestBotID,
		GuestBotName: vote.GuestBotName,
		GuestBotUser: vote.GuestBotUserName,
		Reason:       ReasonVote,
	})
	if err != nil {
		guestSpamTelegram.AnswerCallback(callback.ID, true, "投票通过，已拉黑，但删除消息失败，请管理员检查权限")
		guestSpamTelegram.DeleteCallbackMessage(callback)
		return err
	}
	guestSpamTelegram.AnswerCallback(callback.ID, false, "投票通过，已拉黑并删除消息")
	guestSpamTelegram.DeleteCallbackMessage(callback)
	return nil
}
```

并在 `SpamVoteCallback` 中接住返回值：

```go
if vote.VoteScore >= vote.RequiredVoteCount {
	return applyVotePassed(vote, callback)
}
```

- [ ] **Step 4: Run the deletion-focused tests to verify the new callback and logging behavior**

Run:

```powershell
go test ./plugins/antispam -run "TestDeleteGuestMessageLogsFailure|TestApplyVotePassedShowsDeleteFailureMessage|TestCommandAndCallbackGuards" -count=1
go test -tags=integration ./plugins/antispam -run "TestGuestSpamIntegrationHandleLowTrustDeleteFailureStillBlacklists|TestGuestSpamIntegrationVotePassDeleteFailureIsLoggedAndReported|TestGuestSpamIntegrationSpamVoteCallbackPaths" -count=1
```

Expected:

- 删除 helper 成功和失败路径都记录正确日志
- 投票通过但删消息失败时，回调文案明确说“已拉黑，但删除消息失败”
- 低信誉自动删除路径保持兼容

- [ ] **Step 5: Commit the deletion handling fix**

```bash
git add src/plugins/antispam/message.go src/plugins/antispam/commands.go src/plugins/antispam/guest_bot_spam_test.go src/plugins/antispam/guest_bot_spam_integration_test.go
git commit -m "fix: surface guest spam delete failures"
```

---

### Task 4: Update Integration Tests And Run Full Verification

**Files:**
- Modify: `src/plugins/antispam/guest_bot_spam_integration_test.go`
- Modify: `src/plugins/antispam/guest_bot_spam_test.go`

- [ ] **Step 1: Replace remaining set-based test setup with sorted-set setup**

把现有 integration tests 中所有直接写活跃用户的代码从：

```go
bot.GoRedis.SAdd(redisCtx, activeUsersKey(integrationChatID), 1, 2, 3)
```

统一替换为：

```go
for _, userID := range []int64{1, 2, 3} {
	if err := bot.GoRedis.ZAdd(redisCtx, activeUsersKey(integrationChatID), &redis.Z{
		Score:  float64(time.Now().Unix()),
		Member: fmt.Sprintf("%d", userID),
	}).Err(); err != nil {
		t.Fatalf("seed active user %d: %v", userID, err)
	}
}
```

并把所有 `SAdd` 活跃用户预置点都改掉，至少包括：

- `TestGuestSpamIntegrationGuestSpamHandleCandidatesAndStartVote`
- `TestGuestSpamIntegrationGuestSpamHandleSendFailures`
- `TestGuestSpamIntegrationSelectRecentGuestCallbackPaths`
- `TestGuestSpamIntegrationSpamVoteCallbackPaths`
- `TestGuestSpamIntegrationSpamVoteWeightTiers`

- [ ] **Step 2: Run the package test suite locally**

Run:

```powershell
go test ./plugins/antispam -count=1
```

Expected:

- `arknights_bot/plugins/antispam` PASS

- [ ] **Step 3: Run the full repository unit test suite locally**

Run:

```powershell
go test ./... -count=1
```

Expected:

- 仓库内所有非 integration tests PASS

- [ ] **Step 4: Push the branch and rely on CI for integration coverage**

Run:

```bash
git status --short
git push origin codex/intercept-guest-bot
```

Expected:

- 工作区只包含本次修复
- CI 自动运行：
  - `go test ./...`
  - `go test -tags=integration ./plugins/antispam`

在 CI 结果出来后，需要人工确认：

- active-user 相关 integration tests 通过
- restore 相关 integration tests 通过
- vote pass delete failure 相关 integration tests 通过

- [ ] **Step 5: Commit any final local test-only adjustments if needed**

如果在本地验证阶段对测试辅助函数还有最后一轮调整，使用下面的提交模板：

```bash
git add src/plugins/antispam/guest_bot_spam_test.go src/plugins/antispam/guest_bot_spam_integration_test.go
git commit -m "test: align antispam coverage with sorted set flow"
```

如果没有额外调整，这一步标记完成但不新增提交。

---

## Self-Review

### Spec Coverage

- 活跃用户改用 `sorted set`：Task 1
- restore 先 unban 再 unrestrict：Task 2
- 删除 helper 统一且回调文案真实：Task 3
- integration tests 与 CI 验证：Task 4

无遗漏项。

### Placeholder Scan

已检查本计划，没有 `TODO`、`TBD`、`implement later`、`similar to Task N` 等占位项。

### Type Consistency

本计划中使用的函数名保持一致：

- `activeUsersKey`
- `pruneActiveUsers`
- `deleteGuestMessageWithLog`
- `deleteGuestMessageByFields`
- `UnbanChatMember`
- `applyVotePassed`

后续实现时如命名有调整，需要同步修改所有任务中的代码片段与命令。
