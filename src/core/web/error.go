package web

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"html"
	"net/http"
)

// renderError 渲染错误页，Screenshot 检测到错误标记后会把文本作为截图失败原因返回
func renderError(c *gin.Context, err error) {
	body := fmt.Sprintf(`<!DOCTYPE html><html><head><meta charset="utf-8"></head><body><div id="main" class="error" style="position:absolute;width:700px;background-color:#2e3031;color:#fff;font-size:24px;line-height:1.8;padding:40px;font-family:'NotoSansHans',sans-serif;">%s</div></body></html>`, html.EscapeString(err.Error()))
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(body))
}
