# Antispam Guest Spam Fixes Design

## Goal

修复 `antispam` 插件中已确认的三个行为缺陷：

1. 活跃用户统计因为共享 TTL 的 Redis `set` 被长期放大。
2. `/guest_spam_log restore` 只能解除禁言，不能恢复被封禁用户。
3. 投票通过后的删消息失败会被错误地报告为成功。

本次改动只修正这三个问题及其直接测试覆盖，不引入新的产品功能，也不调整数据库表结构。

## Scope

本设计覆盖以下代码路径：

- `src/plugins/antispam/cache.go`
- `src/plugins/antispam/keys.go`
- `src/plugins/antispam/commands.go`
- `src/plugins/antispam/message.go`
- `src/plugins/antispam/telegram_gateway.go`
- `src/plugins/antispam/telegram_gateway_test.go`
- `src/plugins/antispam/guest_bot_spam_test.go`
- `src/plugins/antispam/guest_bot_spam_integration_test.go`

不在本次范围内的内容：

- 变更 `guest_spam_*` 表结构
- 调整投票权重或投票门槛规则
- 引入新的管理命令
- 对现有 `gatekeeper` 插件做联动改造

## Current Problems

### 1. Active User Window Is Incorrect

当前“最近 10 分钟活跃用户”基于 `activeUsersKey(chatID)` 的 Redis `set` 实现，并对整个 key 续期。只要任意用户继续发言，整个集合都会被续期，导致已经超过 10 分钟未发言的老用户依然留在集合中。

这会造成两个直接后果：

- `ActiveUserCount` 会长期偏大，投票门槛被压低。
- `IsActiveUser` 会对实际上已经过期的成员继续返回 `true`。

### 2. Restore Path Does Not Undo Ban

`restoreCaller` 目前只调用 `RestrictChatMember(..., AllPermissions)`。该调用能够解除禁言，但不能撤销通过 `BanChatMember` 施加的封禁。因此管理员看到“已恢复”时，用户可能仍然无法重新加入或发言。

### 3. Vote-Pass Delete Failure Is Silent

投票通过路径在调用 `DeleteMessage` 后忽略错误，并固定回调“已拉黑并删除消息”。这会让管理员和投票参与者误判消息已经删除，同时缺失故障日志，和低信誉自动删除路径的行为也不一致。

## Chosen Approach

### Active User Tracking

采用用户确认的方案 B：将活跃窗口改为 Redis `sorted set`。

每个群使用一个 `sorted set` 保存候选活跃用户：

- key: `guest_spam:active:<chatID>`
- member: `<userID>`
- score: 最近一次普通发言的 Unix 时间戳

写入规则：

- `RecordMessageActivity` 在记录正常发言时执行 `ZADD activeUsersKey(chatID) score userID`
- 写入完成后立即执行一次过期成员清理，删除 `score < now - activeWindowTTL` 的成员
- key 本身仍可设置一个较短 TTL 作为缓存兜底，但逻辑有效期只由 score 决定

读取规则：

- `ActiveUserCount` 调用前先清理过期成员，再读取 `ZCARD`
- `IsActiveUser` 调用前先清理过期成员，再检查成员是否仍存在于 `sorted set`
- `LoadCacheFromDB` 在重建缓存时，使用 `risk.LastMessageAt` 作为 score，只为仍处于活跃窗口内的成员回填 `sorted set`

选择 `sorted set` 的原因：

- 状态表达与“按时间窗口保活”问题直接匹配。
- 避免共享 TTL 带来的整体续命问题。
- 查询活跃人数和单用户活跃状态都可以在现有 Redis 结构中完成，无需额外按用户散 key。
- 这个功能尚未正式上线，当前阶段做结构性修正比最小修补更稳妥。

### Restore Workflow

`restoreCaller` 改为“全恢复”语义，而不是“仅解除禁言”：

1. 调用 `UnbanChatMember(chatID, userID)`
2. 调用 `RestrictChatMember(chatID, userID, tgbotapi.AllPermissions)`
3. 两步均成功后再执行 `RestoreCallerState`
4. 任一步失败都返回错误并记录日志，不发送成功提示

这样可以统一覆盖两类历史状态：

- 用户之前被 `BanChatMember`
- 用户之前只被 `RestrictChatMember`

不额外判断当前到底是 ban 还是 mute，直接按幂等恢复流程处理，降低状态分支复杂度。

### Unified Delete Result Handling

将删除 guest spam 消息的动作统一收敛到一个 helper，供以下路径共享：

- 低信誉自动拦截后的消息删除
- 投票通过后的消息删除

helper 的职责：

- 调用 Telegram 删除消息
- 成功时记录 `ActionDeleteMessage`
- 失败时记录 `ActionDeleteFailed`
- 将错误返回给调用方

投票通过路径基于 helper 结果调整回调文案：

- 删除成功：`投票通过，已拉黑并删除消息`
- 删除失败：`投票通过，已拉黑，但删除消息失败，请管理员检查权限`

这样可以让日志、回调文案和真实 Telegram 状态保持一致。

## Detailed Design

### 1. Redis Key Changes

`src/plugins/antispam/keys.go` 中保留活跃 key 的命名入口，但其语义从“共享 TTL 的活跃集合”改为“按时间戳维护的活跃排序集合”。

新增或调整的辅助函数应满足：

- 能返回当前时间窗口的截止时间戳
- 能执行活跃集合的过期成员清理
- 便于 `RecordMessageActivity`、`ActiveUserCount`、`IsActiveUser`、`LoadCacheFromDB` 复用

本次不新增数据库字段，因为现有 `MemberRisk.LastMessageAt` 已足够支持冷启动回填。

### 2. Cache Layer Updates

`src/plugins/antispam/cache.go` 需要做三类调整：

1. 写路径
   - `RecordMessageActivity` 改写活跃集合更新逻辑
   - 清理过期成员后再更新 Redis 中的活跃人数状态

2. 读路径
   - `ActiveUserCount`
   - `IsActiveUser`
   - 两者都必须在读取前触发过期清理，避免读到陈旧成员

3. 重载路径
   - `LoadCacheFromDB` 在加载 `MemberRisk` 后，仅对 `LastMessageAt` 仍位于窗口内的成员回填到 `sorted set`

这样做之后，活跃成员资格由“最近一次普通发言时间是否落在窗口内”决定，而不是依赖任意群成员继续发言。

### 3. Telegram Gateway Updates

`src/plugins/antispam/telegram_gateway.go` 的接口扩展为支持解除封禁：

- 新增 `UnbanChatMember(chatID, userID int64)` 方法

对应的 fake gateway 需要补齐：

- 调用记录
- 可注入错误
- 测试辅助断言所需字段

这样 `commands.go` 不需要直接依赖 `bot.Arknights`，继续维持动作层的可替换性。

### 4. Restore Path Updates

`src/plugins/antispam/commands.go` 中的 `restoreCaller` 需要显式串行执行恢复动作：

1. `UnbanChatMember`
2. `RestrictChatMember(..., AllPermissions)`
3. `RestoreCallerState`
4. 成功提示

错误处理要求：

- 任一步 Telegram 调用失败都必须 `AddLog`
- 失败时不得清除 warning 或 mute level
- 失败时不得发送“已恢复”类消息

保留当前 `/guest_spam_log restore <userID>` 命令格式，不改 CLI 交互。

### 5. Delete Handling Updates

`src/plugins/antispam/message.go` 和 `src/plugins/antispam/commands.go` 共享新的删除 helper。

低信誉路径：

- `GuestBotSpamHandle` 中删除消息继续执行，但由统一 helper 负责日志和错误返回

投票通过路径：

- `applyVotePassed` 在调用 `ApplyVotePassedState` 后使用统一 helper
- 需要根据 helper 返回结果决定 callback 文案
- 即使消息删除失败，bot 黑名单和投票状态更新仍然保留成功语义

这里的原则是：

- 黑名单判定成功和消息删除成功是两个独立结果
- 不能因为删消息失败就回滚投票通过后的黑名单状态

## Error Handling

### Active User Cleanup

Redis 清理动作失败时：

- `ActiveUserCount` 返回保守结果
- `IsActiveUser` 返回保守结果
- 日志中记录 Redis 读写错误

这里的“保守结果”指遵循现有代码风格，不因为 guest spam 辅助逻辑直接中断整个消息处理链。

### Restore Failures

恢复用户时：

- `UnbanChatMember` 失败，立即记录日志并返回
- `RestrictChatMember` 失败，立即记录日志并返回
- 不做状态清理

这样可避免“本地状态被恢复，但 Telegram 权限状态未恢复”的假成功。

### Delete Failures

删除失败时：

- 统一写 `ActionDeleteFailed`
- 保留错误细节到 `SpamLog.Detail`
- 投票通过路径回调提示管理员检查权限

## Testing Strategy

### Unit Tests

更新或新增以下单测覆盖：

- 活跃用户在 `sorted set` 中按时间窗口失效
- `ActiveUserCount` 读取前会清理过期成员
- `IsActiveUser` 对过期成员返回 `false`
- `restoreCaller` 会先 `UnbanChatMember` 再解除限制
- `restoreCaller` 任一步失败都不会清理 warning
- 投票通过删除失败时，回调文案不再声称“已删除”
- 删除 helper 在成功和失败两种情况下都会写正确日志

### Integration Tests

更新或新增 integration 场景：

- 冷启动 `LoadCacheFromDB` 后，仅窗口内成员被恢复为活跃用户
- 投票人数统计不会因为旧成员残留而虚高
- restore 对已 ban 用户能真正恢复
- vote pass 删除失败时日志与回调一致

CI 已经具备 MySQL 和 Redis，因此 integration tests 继续依赖现有 workflow 即可。

## Rollout Notes

该功能尚未正式上线，因此允许在当前阶段修正 Redis 活跃模型。

上线前需要确认：

- 本分支的 CI 包含 `go test ./...`
- 本分支的 CI 包含 `go test -tags=integration ./plugins/antispam`
- 修复完成后再观察是否需要补充运行时日志指标；本次设计不额外引入监控

## Acceptance Criteria

满足以下条件即可视为本次修复完成：

1. 10 分钟前已停止发言的用户不会因为其他人继续发言而保持活跃资格。
2. `ActiveUserCount` 和 `IsActiveUser` 均基于最近一次发言时间窗口判断。
3. `/guest_spam_log restore <userID>` 能恢复被 ban 的用户，也能恢复被 mute 的用户。
4. restore 任一步失败时，不会错误提示成功，也不会清理本地 warning 状态。
5. 投票通过后如果删消息失败，日志会记录失败，用户提示也会明确说明“已拉黑但删除失败”。
6. 现有低信誉自动拦截、投票、日志查询流程保持可用。
