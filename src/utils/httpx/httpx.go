// Package httpx 提供带兜底的 HTTP 抓取：请求出错或状态码异常时返回错误并记录日志，
// 不会像以前那样忽略 error 后在 nil response 上取 Body，直接把 goroutine 打成空指针崩溃。
package httpx

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
)

// Get 抓取 URL 正文并校验状态码
func Get(rawURL string) ([]byte, error) {
	response, err := http.Get(rawURL)
	if err != nil {
		log.Println("请求失败:", rawURL, err)
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("请求失败，状态码：%d，地址：%s", response.StatusCode, rawURL)
		return nil, fmt.Errorf("状态码 %d：%s", response.StatusCode, rawURL)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("读取正文失败:", rawURL, err)
		return nil, err
	}
	return body, nil
}

// Body 返回正文读取器，抓取失败时返回空正文（调用方解析到空内容而不是 panic）
func Body(rawURL string) io.Reader {
	body, err := Get(rawURL)
	if err != nil {
		return bytes.NewReader(nil)
	}
	return bytes.NewReader(body)
}
