package skland

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/resty/v2"
	"log"
)

const (
	// 服务器
	serverCN     = "国服"
	serverGlobal = "国际服"

	// 森空岛
	sklandAddr = "https://zonai.skland.com"
	skportAddr = "https://zonai.skport.com"

	// Hypergryph
	hypergryphAddr   = "https://as.hypergryph.com"
	hypergryphAKAddr = "https://ak.hypergryph.com"
	gryphlineAddr    = "https://as.gryphline.com"
)

// SKBaseResp 森空岛接口通用响应
type SKBaseResp[T any] struct {
	Code    *int   `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func (t *SKBaseResp[T]) IsSuccess() bool {
	return t != nil && t.Code != nil && *t.Code == 0
}

func (t *SKBaseResp[T]) String() string {
	if t != nil && t.Code != nil {
		return fmt.Sprintf("code: %d, message: %s", *t.Code, t.Message)
	}
	return "<nil>"
}

// HBaseResp Hypergryph 接口通用响应
type HBaseResp[T any] struct {
	StatusCode *int   `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`

	Status *int   `json:"status"`
	Type   string `json:"type"`
	Msg    string `json:"msg"`

	Data T `json:"data,omitempty"`
}

func (t *HBaseResp[T]) IsSuccess() bool {
	return t != nil && t.Status != nil && *t.Status == 0
}

func (t *HBaseResp[T]) String() string {
	if t != nil && t.StatusCode != nil {
		return fmt.Sprintf("status: %d, error: %s, message: %s", *t.StatusCode, t.Error, t.Message)
	} else if t != nil && t.Status != nil {
		return fmt.Sprintf("status: %d, type: %s, msg: %s", *t.Status, t.Type, t.Msg)
	}
	return "<nil>"
}

// SKR 创建森空岛请求
func SKR() *resty.Request {
	r := resty.New()
	if proxy := viper.GetString("proxy"); proxy != "" {
		r.SetProxy(proxy)
	}
	return r.R().
		SetHeader("User-Agent", "Skland/1.21.0 (com.hypergryph.skland; build:102100065; iOS 17.6.0; ) Alamofire/5.7.1").
		SetHeader("Accept-Encoding", "gzip").
		SetHeader("Connection", "close").
		SetHeader("Content-Type", "application/json")
}

// HR 创建 Hypergryph 请求
func HR() *resty.Request {
	return resty.R().
		SetHeader("User-Agent", viper.GetString("api.user_agent")).
		SetHeader("Accept-Encoding", "gzip")
}

// skRequest 森空岛通用请求，isGlobal 为 true 时请求国际服 skport
func skRequest[T any](r *resty.Request, method, path string, isGlobal bool, vs ...any) (t T, _ error) {
	addr, name := sklandAddr, "skland"
	if isGlobal {
		addr, name = skportAddr, "skport"
	}
	for i := 0; i < len(vs); i++ {
		switch v := vs[i].(type) {
		case AccountSkland:
			addSign(r, method, path, v)
		}
	}
	resp, respErr := r.SetError(&SKBaseResp[any]{}).SetResult(&SKBaseResp[T]{}).Execute(method, addr+path)
	if resp.StatusCode() == 405 {
		log.Println(string(resp.Body()))
		return t, fmt.Errorf("服务器被墙了！")
	}
	if resp.StatusCode() == 401 {
		log.Println(string(resp.Body()))
		return t, fmt.Errorf("cred无效！")
	}
	res, err := resty.ParseResp[*SKBaseResp[any], *SKBaseResp[T]](
		resp, respErr,
	)
	if err != nil {
		return t, fmt.Errorf("[%s] %w", name, err)
	}
	return res.Data, nil
}

// skRequestData 森空岛请求，返回原始响应体
func skRequestData(r *resty.Request, method, path string, isGlobal bool, vs ...any) (string, error) {
	addr, name := sklandAddr, "skland"
	if isGlobal {
		addr, name = skportAddr, "skport"
	}
	for i := 0; i < len(vs); i++ {
		switch v := vs[i].(type) {
		case AccountSkland:
			addSign(r, method, path, v)
		}
	}
	res, err := r.Execute(method, addr+path)
	if err != nil {
		return "", fmt.Errorf("[%s] %w", name, err)
	}
	if res.StatusCode() == 405 {
		log.Println(string(res.Body()))
		return "", fmt.Errorf("服务器被墙了！")
	}
	return string(res.Body()), nil
}

// hgRequest Hypergryph 通用请求，isGlobal 为 true 时请求国际服 gryphline
func hgRequest[T any](r *resty.Request, method, path string, isGlobal bool) (t T, _ error) {
	addr, name := hypergryphAddr, "hypergryph"
	if isGlobal {
		addr, name = gryphlineAddr, "gryphline"
	}
	res, err := resty.ParseResp[*HBaseResp[any], *HBaseResp[T]](
		r.SetError(&HBaseResp[any]{}).SetResult(&HBaseResp[T]{}).Execute(method, addr+path),
	)
	if err != nil {
		return t, fmt.Errorf("[%s] %w", name, err)
	}
	return res.Data, nil
}

// hgRawRequest Hypergryph 请求，返回原始响应体
func hgRawRequest(r *resty.Request, method, path, addr string) (string, error) {
	res, err := r.Execute(method, addr+path)
	if err != nil {
		return "", fmt.Errorf("[hypergryph] %w", err)
	}
	return string(res.Body()), nil
}
