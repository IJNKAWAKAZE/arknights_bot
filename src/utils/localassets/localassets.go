// Package localassets 把图床素材缓存到本地磁盘。
//
// 开启后 /update 会把干员的头像、半身像、立绘、皮肤等素材增量下载到本地，
// 渲染图片时优先使用本地文件，本地缺失或下载失败时再回退到远程图床，
// 避免图床临时不可用或网络抖动导致素材加载失败（截图中出现占位图）。
package localassets

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

const (
	// PathPrefix 本地素材在 Web 服务上对应的访问前缀
	PathPrefix = "/local-assets/"
	// DefaultDir 默认的本地素材目录
	DefaultDir = "./assets_local"

	// missingLogLimit 404 这类「图床确实没有」的失败最多单独打几条日志，其余只计入汇总
	missingLogLimit = 5

	defaultConcurrency = 4
	defaultRetries     = 3
	requestTimeout     = 60 * time.Second
)

// Hosts 需要优先使用本地缓存的图床域名，其他域名（例如 wiki 页面链接）保持原样
var Hosts = []string{"media.prts.wiki", "web.hycdn.cn"}

var (
	httpClient = &http.Client{Timeout: requestTimeout}

	// ErrNotFound 图床返回 404：素材确实不存在（未上架的皮肤、没有立绘的召唤物等），重试没有意义
	ErrNotFound = errors.New("素材未提供（404）")
)

// Enabled 是否启用本地素材缓存
func Enabled() bool {
	return viper.GetBool("local_assets.enable")
}

// Dir 本地素材目录
func Dir() string {
	dir := strings.TrimSpace(viper.GetString("local_assets.dir"))
	if dir == "" {
		return DefaultDir
	}
	return dir
}

// RemoteFallback 本地缺失且补齐失败时，是否回退到远程图床
func RemoteFallback() bool {
	if !viper.IsSet("local_assets.remote_fallback") {
		return true
	}
	return viper.GetBool("local_assets.remote_fallback")
}

// Concurrency 素材下载并发数
func Concurrency() int {
	concurrency := viper.GetInt("local_assets.concurrency")
	if concurrency <= 0 {
		concurrency = defaultConcurrency
	}
	if concurrency > 32 {
		concurrency = 32
	}
	return concurrency
}

// Retries 单张素材的最大尝试次数
func Retries() int {
	retries := viper.GetInt("local_assets.retry")
	if retries <= 0 {
		retries = defaultRetries
	}
	return retries
}

// IsRemoteHost 是否是走本地缓存的图床域名（比较时忽略端口）
func IsRemoteHost(host string) bool {
	h := normalizeHost(host)
	for _, item := range Hosts {
		if h == normalizeHost(item) {
			return true
		}
	}
	return false
}

// normalizeHost 归一化域名，比较时忽略大小写与端口
func normalizeHost(host string) string {
	h := strings.ToLower(strings.TrimSpace(host))
	if i := strings.LastIndex(h, ":"); i >= 0 {
		h = h[:i]
	}
	return h
}

// EscapePath 按图床（MediaWiki）的规则对文件名或路径做百分号编码：
// +、#、空格、中文这些字符必须转义，否则 media.prts.wiki 会把
// 「立绘_阿米娅_1+.png」这类地址当成另一个文件，直接返回 404。路径分隔符 / 保留。
func EscapePath(p string) string {
	const hex = "0123456789ABCDEF"
	var b strings.Builder
	b.Grow(len(p) * 3 / 2)
	for i := 0; i < len(p); i++ {
		c := p[i]
		switch {
		case c == '/' || c == '-' || c == '_' || c == '.' || c == '~',
			c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9',
			// MediaWiki 也不转义这几个字符，保持地址可读
			c == '(', c == ')', c == '!', c == '*', c == '\'':
			b.WriteByte(c)
		default:
			b.WriteByte('%')
			b.WriteByte(hex[c>>4])
			b.WriteByte(hex[c&0x0f])
		}
	}
	return b.String()
}

// EscapedURL 把地址的路径重新编码成图床能识别的形式。
// 请求图床、以及在页面上回退远程地址时都用它，保证同一个素材只会对应一份本地文件。
func EscapedURL(rawURL string) string {
	u, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || u.Host == "" || u.Path == "" {
		return rawURL
	}
	u.RawPath = EscapePath(u.Path)
	return u.String()
}

// CachePath 把远程素材地址映射成本地缓存文件路径：<dir>/<host>/<原始路径>。
// 带 image_process=format,webp 的地址由图床转成 webp，因此按 webp 落盘。
func CachePath(rawURL string) (string, bool) {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return "", false
	}
	u, err := url.Parse(rawURL)
	if err != nil || u.Host == "" || !IsRemoteHost(u.Host) {
		return "", false
	}
	clean := path.Clean("/" + strings.Trim(u.Path, "/"))
	if clean == "/" || clean == "." {
		return "", false
	}
	segments := strings.Split(strings.TrimPrefix(clean, "/"), "/")
	for i, segment := range segments {
		segments[i] = sanitizeSegment(segment)
	}
	rel := strings.Join(segments, "/")
	if ext := path.Ext(rel); ext == "" {
		rel += ".png"
	} else if wantsWebp(u) {
		rel = strings.TrimSuffix(rel, ext) + ".webp"
	}
	// 域名做目录名，去掉端口等不适合做文件名的字符
	target := filepath.Join(Dir(), sanitizeSegment(strings.ToLower(u.Host)), filepath.FromSlash(rel))
	if !withinDir(target) {
		return "", false
	}
	return target, true
}

// wantsWebp 地址是否要求图床转成 webp
func wantsWebp(u *url.URL) bool {
	return strings.Contains(strings.ToLower(u.Query().Get("image_process")), "webp")
}

// sanitizeSegment 去掉不适合做本地文件名的字符
func sanitizeSegment(segment string) string {
	segment = strings.Map(func(r rune) rune {
		switch r {
		case '<', '>', ':', '"', '/', '\\', '|', '?', '*':
			return '_'
		}
		if r < 32 {
			return '_'
		}
		return r
	}, segment)
	segment = strings.TrimRight(segment, " .")
	if segment == "" || segment == "." || segment == ".." {
		return "_"
	}
	return segment
}

// withinDir 确保缓存路径没有跑到素材目录之外
func withinDir(target string) bool {
	base, err := filepath.Abs(Dir())
	if err != nil {
		return false
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return false
	}
	rel, err := filepath.Rel(base, abs)
	if err != nil {
		return false
	}
	return rel != ".." && !strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}

// Cached 远程素材是否已经在本地
func Cached(rawURL string) bool {
	target, ok := CachePath(rawURL)
	if !ok {
		return false
	}
	return fileExists(target)
}

// Read 读取本地缓存的素材，返回内容与图片类型
func Read(rawURL string) ([]byte, string, bool) {
	target, ok := CachePath(rawURL)
	if !ok {
		return nil, "", false
	}
	data, err := os.ReadFile(target)
	if err != nil || len(data) == 0 {
		return nil, "", false
	}
	return data, detectContentType(target, data), true
}

// ContentTypeForFile 根据文件头判断图片类型，判断不出来时退回扩展名
func ContentTypeForFile(target string) string {
	file, err := os.Open(target)
	if err != nil {
		return contentTypeByExt(target)
	}
	defer file.Close()
	head := make([]byte, 16)
	n, err := io.ReadFull(file, head)
	if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
		return contentTypeByExt(target)
	}
	return detectContentType(target, head[:n])
}

func detectContentType(name string, head []byte) string {
	switch {
	case len(head) >= 12 && string(head[0:4]) == "RIFF" && string(head[8:12]) == "WEBP":
		return "image/webp"
	case len(head) >= 8 && string(head[0:8]) == "\x89PNG\r\n\x1a\n":
		return "image/png"
	case len(head) >= 3 && head[0] == 0xFF && head[1] == 0xD8 && head[2] == 0xFF:
		return "image/jpeg"
	case len(head) >= 6 && (string(head[0:6]) == "GIF87a" || string(head[0:6]) == "GIF89a"):
		return "image/gif"
	}
	return contentTypeByExt(name)
}

func contentTypeByExt(name string) string {
	switch strings.ToLower(filepath.Ext(name)) {
	case ".webp":
		return "image/webp"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	}
	return "application/octet-stream"
}

func fileExists(target string) bool {
	info, err := os.Stat(target)
	return err == nil && info.Mode().IsRegular() && info.Size() > 0
}

// Download 下载单张素材到本地，已存在时直接跳过；
// 返回新增的字节数以及本次是否真的下载了文件。
func Download(rawURL string) (int64, bool, error) {
	target, ok := CachePath(rawURL)
	if !ok {
		return 0, false, fmt.Errorf("不支持的素材地址：%s", rawURL)
	}
	if fileExists(target) {
		return 0, false, nil
	}
	// 同一张图片可能被多个请求同时发现缺失，这里加锁避免重复下载
	lock := keyLock(rawURL)
	lock.Lock()
	defer lock.Unlock()
	if fileExists(target) {
		return 0, false, nil
	}
	size, err := fetchToFile(rawURL, target)
	if err != nil {
		return 0, false, err
	}
	return size, true, nil
}

// fetchToFile 下载并原子落盘，失败按配置自动重试
func fetchToFile(rawURL, target string) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return 0, err
	}
	attempts := Retries()
	var lastErr error
	for attempt := 1; attempt <= attempts; attempt++ {
		size, retryable, err := tryFetch(rawURL, target)
		if err == nil {
			return size, nil
		}
		lastErr = err
		if !retryable {
			break
		}
		if attempt < attempts {
			time.Sleep(time.Duration(attempt) * time.Second)
		}
	}
	return 0, lastErr
}

// tryFetch 下载单个文件，并返回这个错误是否值得重试
func tryFetch(rawURL, target string) (int64, bool, error) {
	// 文件名里的 +、# 等字符不转义会被图床当成另一个文件（404），这里统一编码
	resp, err := httpClient.Get(EscapedURL(rawURL))
	if err != nil {
		return 0, true, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		// 404 之类的错误重试也没用，免得白等、白占图床带宽
		retryable := resp.StatusCode >= 500 ||
			resp.StatusCode == http.StatusRequestTimeout ||
			resp.StatusCode == http.StatusTooManyRequests
		if resp.StatusCode == http.StatusNotFound {
			return 0, false, fmt.Errorf("%w：%s", ErrNotFound, rawURL)
		}
		return 0, retryable, fmt.Errorf("状态码 %d", resp.StatusCode)
	}
	tmp := target + ".download"
	file, err := os.Create(tmp)
	if err != nil {
		return 0, true, err
	}
	size, copyErr := io.Copy(file, resp.Body)
	closeErr := file.Close()
	if copyErr == nil {
		copyErr = closeErr
	}
	if copyErr != nil {
		os.Remove(tmp)
		return 0, true, copyErr
	}
	if size == 0 {
		os.Remove(tmp)
		return 0, true, errors.New("图片内容为空")
	}
	// 先写临时文件再改名，避免半截文件被当成缓存的图片
	if err := os.Rename(tmp, target); err != nil {
		os.Remove(tmp)
		return 0, true, err
	}
	return size, false, nil
}

// Result 一次素材同步的结果
type Result struct {
	Total      int
	Downloaded int
	Skipped    int
	Failed     int
	// Missing 其中「图床确实没有素材」（404）的数量，属于预期内，不算真正的故障
	Missing  int
	Bytes    int64
	Duration time.Duration
}

// Summary 同步结果的简要描述
func (r Result) Summary() string {
	text := fmt.Sprintf("本地素材同步：共 %d 张，新增 %d 张（%s），已存在 %d 张，失败 %d 张，耗时 %s",
		r.Total, r.Downloaded, humanSize(r.Bytes), r.Skipped, r.Failed, r.Duration.Round(time.Second))
	if r.Missing > 0 {
		text += fmt.Sprintf("（其中 %d 张图床未提供，属预期内：未上架皮肤、没有立绘的召唤物等）", r.Missing)
	}
	return text
}

func humanSize(size int64) string {
	switch {
	case size >= 1<<30:
		return fmt.Sprintf("%.2f GB", float64(size)/(1<<30))
	case size >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(size)/(1<<20))
	case size >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(size)/(1<<10))
	}
	return fmt.Sprintf("%d B", size)
}

// Sync 增量同步一批素材，已经在本地的不再重复下载
func Sync(urls []string) Result {
	start := time.Now()
	list := uniqueURLs(urls)
	result := Result{Total: len(list)}
	if len(list) == 0 {
		result.Duration = time.Since(start)
		return result
	}
	log.Printf("本地素材同步开始，待处理 %d 张，并发 %d", len(list), Concurrency())
	concurrency := Concurrency()
	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup
	var mu sync.Mutex
	done := 0
	for _, rawURL := range list {
		wg.Add(1)
		sem <- struct{}{}
		go func(rawURL string) {
			defer wg.Done()
			defer func() { <-sem }()
			size, downloaded, err := Download(rawURL)
			mu.Lock()
			defer mu.Unlock()
			done++
			switch {
			case err != nil:
				result.Failed++
				if errors.Is(err, ErrNotFound) {
					// 图床确实没有这张图，属于预期内，只汇总数量，不逐条刷日志
					result.Missing++
					if result.Missing <= missingLogLimit {
						log.Printf("本地素材未提供，跳过：%s", rawURL)
					}
				} else {
					log.Printf("本地素材下载失败：%s（%v）", rawURL, err)
				}
			case downloaded:
				result.Downloaded++
				result.Bytes += size
			default:
				result.Skipped++
			}
			if done%200 == 0 {
				log.Printf("本地素材同步进度：%d/%d", done, len(list))
			}
		}(rawURL)
	}
	wg.Wait()
	result.Duration = time.Since(start)
	return result
}

func uniqueURLs(urls []string) []string {
	seen := make(map[string]struct{}, len(urls))
	list := make([]string, 0, len(urls))
	for _, rawURL := range urls {
		rawURL = strings.TrimSpace(rawURL)
		if rawURL == "" {
			continue
		}
		if _, ok := seen[rawURL]; ok {
			continue
		}
		seen[rawURL] = struct{}{}
		list = append(list, rawURL)
	}
	return list
}

var (
	syncMu        sync.Mutex
	syncing       bool
	lastMu        sync.Mutex
	lastResult    Result
	hasLastResult bool
)

// SyncAsync 在后台增量同步素材，已有同步任务在执行时直接跳过
func SyncAsync(urls []string) bool {
	if !Enabled() {
		return false
	}
	syncMu.Lock()
	if syncing {
		syncMu.Unlock()
		log.Println("本地素材同步已在进行中，跳过本次")
		return false
	}
	syncing = true
	syncMu.Unlock()
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Println("本地素材同步异常:", err)
			}
			syncMu.Lock()
			syncing = false
			syncMu.Unlock()
		}()
		result := Sync(urls)
		lastMu.Lock()
		lastResult = result
		hasLastResult = true
		lastMu.Unlock()
		log.Println(result.Summary())
	}()
	return true
}

// LastSummary 最近一次素材同步的结果描述，没有同步过时返回空字符串
func LastSummary() string {
	lastMu.Lock()
	defer lastMu.Unlock()
	if !hasLastResult {
		return ""
	}
	return lastResult.Summary()
}

var (
	lockMu   sync.Mutex
	inflight = make(map[string]*sync.Mutex)
)

// keyLock 返回某个地址对应的锁，避免同一张图被并发重复下载
func keyLock(key string) *sync.Mutex {
	lockMu.Lock()
	defer lockMu.Unlock()
	lock, ok := inflight[key]
	if !ok {
		lock = new(sync.Mutex)
		inflight[key] = lock
	}
	return lock
}
