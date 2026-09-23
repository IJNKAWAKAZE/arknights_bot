// Package httpx 提供带超时限制和状态码校验的 HTTP 请求。
package httpx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Client 独立于 http.DefaultClient，仅在启动阶段配置，避免并发修改。
var Client = &http.Client{Timeout: 30 * time.Second}

// Do 返回状态码校验通过的响应，调用方负责关闭响应体。
func Do(request *http.Request) (*http.Response, error) {
	response, err := Client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		response.Body.Close()
		return nil, fmt.Errorf("HTTP状态码 %d", response.StatusCode)
	}
	return response, nil
}
func Open(rawURL string) (*http.Response, error) {
	request, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	return Do(request)
}
func Get(rawURL string) ([]byte, error) {
	response, err := Open(rawURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	return io.ReadAll(response.Body)
}
func Body(rawURL string) io.Reader {
	body, err := Get(rawURL)
	if err != nil {
		log.Println("抓取正文失败:", err)
		return bytes.NewReader(nil)
	}
	return bytes.NewReader(body)
}
