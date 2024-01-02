package utils

import (
	"arknights_bot/config"
	"bytes"
	"context"
	"github.com/go-redis/redis/v8"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"log"
	"strconv"
	"time"
)

var ctx = context.Background()

// GetFullName 获取用户全名
func GetFullName(user *tgbotapi.User) string {
	var buffer bytes.Buffer
	firstName := user.FirstName
	lastName := user.LastName
	if firstName != "" {
		buffer.WriteString(firstName)
	}
	if lastName != "" {
		buffer.WriteString(lastName)
	}
	return buffer.String()
}

type GroupInvite struct {
	Id           string    `json:"id" gorm:"primaryKey"`
	GroupName    string    `json:"groupName"`
	GroupNumber  string    `json:"groupNumber"`
	UserName     string    `json:"userName"`
	UserNumber   string    `json:"userNumber"`
	MemberName   string    `json:"memberName"`
	MemberNumber string    `json:"memberNumber"`
	CreateTime   time.Time `json:"createTime" gorm:"autoCreateTime"`
	UpdateTime   time.Time `json:"updateTime" gorm:"autoUpdateTime"`
	Remark       string    `json:"remark"`
	Deleted      int64     `json:"deleted"`
}

// SaveInvite 保存邀请记录
func SaveInvite(message *tgbotapi.Message, member *tgbotapi.User) {
	id, _ := gonanoid.New(32)
	groupMessage := GroupInvite{
		Id:           id,
		GroupName:    message.Chat.Title,
		GroupNumber:  strconv.FormatInt(message.Chat.ID, 10),
		UserName:     GetFullName(message.From),
		UserNumber:   strconv.FormatInt(message.From.ID, 10),
		MemberName:   GetFullName(member),
		MemberNumber: strconv.FormatInt(member.ID, 10),
		Deleted:      0,
	}

	config.DBEngine.Table("group_invite").Create(&groupMessage)
}

// GetJoinedGroups 获取加入的群组
func GetJoinedGroups() []string {
	var groups []string
	config.DBEngine.Raw("select group_number from group_message where chat_type = 'supergroup' group by group_number").Scan(&groups)
	return groups
}

// RedisSet redis存值
func RedisSet(key string, val interface{}, expiration time.Duration) {
	err := config.GoRedis.Set(ctx, key, val, expiration).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisGet redis取值
func RedisGet(key string) string {
	val, err := config.GoRedis.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return ""
		}
		log.Println(err)
	}
	return val
}

// RedisIsExists 判断redis值是否存在
func RedisIsExists(key string) bool {
	val := RedisGet(key)
	if val == "" {
		return false
	}
	return true
}

// RedisDel redis根据key删除
func RedisDel(key string) {
	err := config.GoRedis.Del(ctx, key).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisScanKeys 扫描匹配keys
func RedisScanKeys(match string) (*redis.ScanIterator, context.Context) {
	return config.GoRedis.Scan(ctx, 0, match, 0).Iterator(), ctx
}

// RedisSetList redis添加链表元素
func RedisSetList(key string, val interface{}) {
	err := config.GoRedis.RPush(ctx, key, val).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisGetList redis获取所有链表元素
func RedisGetList(key string) []string {
	val, err := config.GoRedis.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		log.Println(err)
	}
	return val
}

// RedisDelListItem redis移除链表元素
func RedisDelListItem(key string, val string) {
	err := config.GoRedis.LRem(ctx, key, 0, val).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisAddSet redis集合添加元素
func RedisAddSet(key string, val string) {
	err := config.GoRedis.SAdd(ctx, key, val).Err()
	if err != nil {
		log.Println(err)
	}
}

// RedisSetIsExists redis集合是否包含元素
func RedisSetIsExists(key string, val string) bool {
	exists, err := config.GoRedis.SIsMember(ctx, key, val).Result()
	if err != nil {
		log.Println(err)
	}
	return exists
}

// RedisDelSetItem redis移除集合元素
func RedisDelSetItem(key string, val string) {
	err := config.GoRedis.SRem(ctx, key, val).Err()
	if err != nil {
		log.Println(err)
	}
}
