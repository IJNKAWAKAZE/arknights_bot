package datasource

import (
	"arknights_bot/utils/cache"
	"arknights_bot/utils/model"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/starudream/go-lib/core/v2/codec/json"
)

const (
	// birthdayKeyPrefix 生日缓存的 key 前缀，值为该日期下的干员列表
	birthdayKeyPrefix = "birthday:"
	// unknownBirthday 没查到生日日期的干员统一放在这个分组，下次更新时重新查
	unknownBirthday = "未知"
	// maxRedirectHops 干员页面重定向的最大跟随层数
	maxRedirectHops = 2
)

var (
	birthdayPattern     = regexp.MustCompile(`\|\s*生日\s*=\s*([^\n|}]+)`)
	birthdayDatePattern = regexp.MustCompile(`^\d+月\s*\d+日$`)
	wikiLinkPattern     = regexp.MustCompile(`\[\[([^\]|]+)`)
)

// loadBirthdayCache 读取已有的生日缓存，返回各日期分组以及每个干员当前所在的分组
func loadBirthdayCache() (map[string][]model.Operator, map[string]string) {
	birthdays := make(map[string][]model.Operator)
	buckets := make(map[string]string)
	iterator, ctx := cache.RedisScanKeys(birthdayKeyPrefix + "*")
	for iterator.Next(ctx) {
		key := iterator.Val()
		var list []model.Operator
		if err := json.Unmarshal([]byte(cache.RedisGet(key)), &list); err != nil {
			log.Println("读取生日缓存失败:", key, err)
			continue
		}
		bucket := strings.TrimPrefix(key, birthdayKeyPrefix)
		birthdays[bucket] = list
		for _, operator := range list {
			buckets[operator.Name] = bucket
		}
	}
	return birthdays, buckets
}

// removeFromBirthday 把干员从旧的生日分组里摘掉，避免同一个干员挂在两个日期下
func removeFromBirthday(birthdays map[string][]model.Operator, bucket, name string) {
	list, ok := birthdays[bucket]
	if !ok {
		return
	}
	kept := make([]model.Operator, 0, len(list))
	for _, operator := range list {
		if operator.Name != name {
			kept = append(kept, operator)
		}
	}
	birthdays[bucket] = kept
}

// fetchBirthday 查询干员生日，返回所属分组：真实日期形如 12月23日，
// 页面上写的是「未公开」「本人表示遗忘」这类文本时按原文分组，字段为空或缺失时归入未知分组。
// 直接解析页面原始 wikitext 里的「|生日=」字段，比在渲染后的档案文本里按行号取更稳，
// 页面是重定向时（改名、带后缀的干员页）会跟着重定向再取一次。
// 返回 false 表示这次没查到（网络或状态码异常），缓存里保持原样，下次更新重试。
func fetchBirthday(api, name string) (string, bool) {
	body, ok := rawWikiText(api, name, maxRedirectHops)
	if !ok {
		return "", false
	}
	return parseBirthday(body), true
}

// rawWikiText 取页面原始 wikitext，页面为重定向时跟随
func rawWikiText(api, name string, hops int) ([]byte, bool) {
	body, ok := getRawWikiText(api, name)
	if !ok {
		return nil, false
	}
	target, isRedirect := redirectTarget(body)
	if !isRedirect {
		return body, true
	}
	if hops <= 0 {
		log.Println("干员页面重定向层数过多:", name)
		return nil, false
	}
	return rawWikiText(api, target, hops-1)
}

func getRawWikiText(api, name string) ([]byte, bool) {
	response, err := http.Get(api + name + "?action=raw")
	if err != nil {
		log.Println("获取干员页面失败:", name, err)
		return nil, false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("获取干员页面失败，状态码：%d，干员：%s", response.StatusCode, name)
		return nil, false
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("读取干员页面失败:", name, err)
		return nil, false
	}
	return body, true
}

// parseBirthday 解析生日分组。
// 日期直接按日期分组；页面上写的「未公开」「未录入」「本人表示遗忘」这类文本按原文分组，
// 免得把这些情况都混成一个“未知”；字段为空或缺失才算未知。
func parseBirthday(wikitext []byte) string {
	match := birthdayPattern.FindSubmatch(wikitext)
	if match == nil {
		return unknownBirthday
	}
	birthday := strings.TrimSpace(string(match[1]))
	if index := strings.Index(birthday, "<"); index >= 0 {
		birthday = strings.TrimSpace(birthday[:index])
	}
	if birthday == "" {
		return unknownBirthday
	}
	return birthday
}

// isBirthdayDate 判断分组名是不是真实日期（形如 12月23日）
func isBirthdayDate(bucket string) bool {
	return birthdayDatePattern.MatchString(bucket)
}

// redirectTarget 从重定向页面里取出目标页面名
func redirectTarget(wikitext []byte) (string, bool) {
	text := strings.TrimSpace(string(wikitext))
	if !strings.HasPrefix(strings.ToLower(text), "#redirect") && !strings.HasPrefix(text, "#重定向") {
		return "", false
	}
	match := wikiLinkPattern.FindStringSubmatch(text)
	if match == nil {
		return "", false
	}
	return strings.TrimSpace(match[1]), true
}
