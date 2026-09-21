package localassets

import (
	"bytes"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Register 注册本地素材的改写中间件与访问路由。
//
// 页面里的图床地址会被改写成同源的 /local-assets/<host>/<path>：
// 命中本地缓存直接返回，本地缺失时先尝试补齐，仍失败再按配置回退到远程地址。
// 配置项在每次请求时读取，因此修改开关后重启或热更新配置都能生效。
func Register(r *gin.Engine) {
	r.Use(rewriteMiddleware())
	r.GET("/local-assets/*filepath", assetHandler)
	r.HEAD("/local-assets/*filepath", assetHandler)
	if Enabled() {
		log.Printf("已开启本地素材缓存，目录 %s，本地缺失时%s", Dir(), fallbackDesc())
	}
}

func fallbackDesc() string {
	if RemoteFallback() {
		return "回退远程图床"
	}
	return "显示占位图"
}

// rewriteMiddleware 把 HTML 响应里的图床地址改写成同源地址
func rewriteMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !Enabled() || !needsRewrite(c.Request.URL.Path) {
			c.Next()
			return
		}
		writer := &rewriteWriter{ResponseWriter: c.Writer}
		c.Writer = writer
		c.Next()
		writer.flush()
	}
}

// needsRewrite 静态资源与本地素材本身不需要改写，避免把大图也缓冲进内存
func needsRewrite(path string) bool {
	if strings.HasPrefix(path, PathPrefix) {
		return false
	}
	if strings.HasPrefix(path, "/assets/") || strings.HasPrefix(path, "/template/") {
		return false
	}
	return true
}

// assetHandler 返回本地缓存的素材，缺失时先补齐再回退远程
func assetHandler(c *gin.Context) {
	if !Enabled() {
		c.Status(http.StatusNotFound)
		return
	}
	remote, ok := remoteURL(c.Param("filepath"), c.Request.URL.RawQuery)
	if !ok {
		c.Status(http.StatusNotFound)
		return
	}
	if target, ok := CachePath(remote); ok && fileExists(target) {
		serveAsset(c, target)
		return
	}
	// 本地还没有这张图，先补齐，避免页面因为图床抖动丢图
	if _, _, err := Download(remote); err == nil {
		if target, ok := CachePath(remote); ok && fileExists(target) {
			serveAsset(c, target)
			return
		}
	}
	if RemoteFallback() {
		c.Redirect(http.StatusFound, EscapedURL(remote))
		return
	}
	c.Status(http.StatusNotFound)
}

func serveAsset(c *gin.Context, target string) {
	c.Header("Content-Type", ContentTypeForFile(target))
	c.Header("Cache-Control", "public, max-age=604800")
	c.File(target)
}

// remoteURL 把 /local-assets/<host>/<path> 还原成原始图床地址
func remoteURL(filepath, rawQuery string) (string, bool) {
	filepath = strings.TrimPrefix(filepath, "/")
	host, rest, found := strings.Cut(filepath, "/")
	if !found || host == "" || rest == "" || !IsRemoteHost(host) {
		return "", false
	}
	u := &url.URL{Scheme: "https", Host: host, Path: "/" + rest, RawQuery: rawQuery}
	return u.String(), true
}

// RewriteHTML 把页面里指向图床的地址改写成同源地址（本地缺失时由路由回退远程）
func RewriteHTML(body []byte) []byte {
	if len(body) == 0 {
		return body
	}
	pairs := make([]string, 0, len(Hosts)*4)
	for _, host := range Hosts {
		pairs = append(pairs,
			"https://"+host+"/", PathPrefix+host+"/",
			"http://"+host+"/", PathPrefix+host+"/",
		)
	}
	return []byte(strings.NewReplacer(pairs...).Replace(string(body)))
}

// rewriteWriter 缓冲 HTML 响应，改写后再写回客户端
type rewriteWriter struct {
	gin.ResponseWriter
	status int
	body   bytes.Buffer
}

func (w *rewriteWriter) WriteHeader(code int) {
	if w.status == 0 {
		w.status = code
	}
}

func (w *rewriteWriter) Write(data []byte) (int, error) {
	return w.body.Write(data)
}

func (w *rewriteWriter) WriteString(s string) (int, error) {
	return w.body.WriteString(s)
}

func (w *rewriteWriter) Written() bool {
	return w.status != 0 || w.body.Len() > 0
}

func (w *rewriteWriter) Size() int {
	return w.body.Len()
}

func (w *rewriteWriter) Status() int {
	if w.status == 0 {
		return http.StatusOK
	}
	return w.status
}

// flush 改写缓冲内容并写回
func (w *rewriteWriter) flush() {
	body := w.body.Bytes()
	if strings.Contains(w.Header().Get("Content-Type"), "text/html") {
		body = RewriteHTML(body)
	}
	header := w.ResponseWriter.Header()
	header.Del("Content-Length")
	header.Set("Content-Length", strconv.Itoa(len(body)))
	w.ResponseWriter.WriteHeader(w.Status())
	if len(body) > 0 {
		_, _ = w.ResponseWriter.Write(body)
	}
}
