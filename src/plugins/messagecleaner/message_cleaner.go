package messagecleaner

import (
	"arknights_bot/config"
	"arknights_bot/plugins/commandoperation"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/ijnkawakaze/telegram-bot-api"
)

const (
	legacyKey     = "msgObjects"
	dueKey        = "msgObjects:due"
	batchSize     = 100
	leaseDuration = 5 * time.Minute
)

var cleanerMu sync.Mutex

type MsgObject struct {
	ChatId       int64     `json:"chatId"`
	MessageId    int64     `json:"messageId"`
	CreateTime   time.Time `json:"createTime"`
	DelTime      float64   `json:"delTime"`
	FunctionHash string    `json:"functionHash"`
}

type deletionQueue struct{ client *redis.Client }

// 原子迁移刚读取的队首任务；Redis 操作失败或其他进程已移走该任务时，
// 不删除旧队列中的数据。
var migrateScript = redis.NewScript(`
if redis.call('LINDEX', KEYS[1], 0) ~= ARGV[1] then return 0 end
if ARGV[2] ~= '' then redis.call('ZADD', KEYS[2], 'NX', ARGV[2], ARGV[1]) end
redis.call('LPOP', KEYS[1])
return 1`)
var claimScript = redis.NewScript(`
local items = redis.call('ZRANGEBYSCORE', KEYS[1], '-inf', ARGV[1], 'LIMIT', 0, 1)
if #items == 0 then return '' end
redis.call('ZADD', KEYS[1], ARGV[2], items[1])
return items[1]`)
var finishScript = redis.NewScript(`
if redis.call('ZSCORE', KEYS[1], ARGV[1]) ~= ARGV[2] then return 0 end
if ARGV[3] == '' then redis.call('ZREM', KEYS[1], ARGV[1])
else redis.call('ZADD', KEYS[1], ARGV[3], ARGV[1]) end
return 1`)

func deadline(msg MsgObject) int64 {
	return msg.CreateTime.Add(time.Duration(msg.DelTime * float64(time.Second))).UnixMilli()
}
func (q deletionQueue) add(ctx context.Context, msg MsgObject) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return q.client.ZAdd(ctx, dueKey, &redis.Z{Score: float64(deadline(msg)), Member: string(data)}).Err()
}
func (q deletionQueue) migrate(ctx context.Context) error {
	for i := 0; i < batchSize; i++ {
		payload, err := q.client.LIndex(ctx, legacyKey, 0).Result()
		if errors.Is(err, redis.Nil) {
			return nil
		}
		if err != nil {
			return err
		}
		var msg MsgObject
		score := ""
		if err := json.Unmarshal([]byte(payload), &msg); err != nil {
			log.Println("丢弃无效的旧删消息任务:", err)
		} else {
			score = strconv.FormatInt(deadline(msg), 10)
		}
		if err := migrateScript.Run(ctx, q.client, []string{legacyKey, dueKey}, payload, score).Err(); err != nil {
			return err
		}
	}
	return nil
}
func (q deletionQueue) claim(ctx context.Context, now time.Time) (string, int64, error) {
	lease := now.Add(leaseDuration).UnixMilli()
	payload, err := claimScript.Run(ctx, q.client, []string{dueKey}, now.UnixMilli(), lease).Text()
	return payload, lease, err
}
func (q deletionQueue) finish(ctx context.Context, payload string, lease int64, retryAt time.Time) error {
	score := ""
	if !retryAt.IsZero() {
		score = strconv.FormatInt(retryAt.UnixMilli(), 10)
	}
	return finishScript.Run(ctx, q.client, []string{dueKey}, payload, lease, score).Err()
}

func retryDelay(err error) time.Duration {
	var apiErr *tgbotapi.Error
	if errors.As(err, &apiErr) && apiErr.RetryAfter > 0 {
		return time.Duration(apiErr.RetryAfter) * time.Second
	}
	return 30 * time.Second
}
func deletionComplete(err error) bool {
	if err == nil {
		return true
	}
	var apiErr *tgbotapi.Error
	return errors.As(err, &apiErr) && apiErr.Code == 400 && strings.Contains(strings.ToLower(apiErr.Message), "message to delete not found")
}

// process 每次只领取一条消息并设置处理期限，避免 Telegram 请求过慢导致
// 整批任务的处理期限到期。未确认完成的任务会保留，重启后仍可重试。
func (q deletionQueue) process(ctx context.Context, send func(MsgObject) error, cleanup func(string)) error {
	if err := q.migrate(ctx); err != nil {
		return err
	}
	for i := 0; i < batchSize; i++ {
		payload, lease, err := q.claim(ctx, time.Now())
		if err != nil {
			return err
		}
		if payload == "" {
			return nil
		}
		var msg MsgObject
		if err := json.Unmarshal([]byte(payload), &msg); err != nil {
			log.Println("丢弃无效删消息任务:", err)
			if err := q.finish(ctx, payload, lease, time.Time{}); err != nil {
				return err
			}
			continue
		}
		sendErr := send(msg)
		if !deletionComplete(sendErr) {
			log.Printf("删除消息 %d/%d 失败: %v", msg.ChatId, msg.MessageId, sendErr)
			if err := q.finish(ctx, payload, lease, time.Now().Add(retryDelay(sendErr))); err != nil {
				return err
			}
			continue
		}
		if err := q.finish(ctx, payload, lease, time.Time{}); err != nil {
			return err
		}
		if msg.FunctionHash != "" && msg.FunctionHash != "None" {
			cleanup(msg.FunctionHash)
		}
	}
	return nil
}

func DelMsg() {
	if !cleanerMu.TryLock() {
		return
	}
	defer cleanerMu.Unlock()
	err := (deletionQueue{config.GoRedis}).process(context.Background(), func(msg MsgObject) error {
		_, err := config.Arknights.Request(tgbotapi.NewDeleteMessage(msg.ChatId, msg.MessageId))
		return err
	}, commandoperation.RemoveCallBack)
	if err != nil {
		log.Println("清理消息失败:", err)
	}
}
func AddDelQueue(chatId int64, messageId int64, delTime float64) {
	AddDelQueueFuncHash(chatId, messageId, delTime, "None")
}
func AddDelQueueFuncHash(chatId int64, messageId int64, delTime float64, hash string) {
	msg := MsgObject{ChatId: chatId, MessageId: messageId, CreateTime: time.Now(), DelTime: delTime, FunctionHash: hash}
	if err := (deletionQueue{config.GoRedis}).add(context.Background(), msg); err != nil {
		log.Println(fmt.Errorf("添加删消息任务: %w", err))
	}
}
