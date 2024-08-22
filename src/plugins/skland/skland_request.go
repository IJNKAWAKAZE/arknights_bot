package skland

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/resty/v2"
	"log"
)

var sklandAddr = "https://zonai.skland.com"

type SKBaseResp[T any] struct {
	Code    *int   `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
}

func (t *SKBaseResp[T]) IsSuccess() bool {
	return t != nil && t.Code != nil && *t.Code == 0
}

func (t *SKBaseResp[T]) String() string {
	return fmt.Sprintf("code: %d, message: %s", *t.Code, t.Message)
}

func SKR() *resty.Request {
	r := resty.New()
	proxy := viper.GetString("proxy")
	if proxy != "" {
		r.SetProxy(proxy)
	}
	return r.R().SetHeader("User-Agent", "Skland/1.21.0 (com.hypergryph.skland; build:102100065; iOS 17.6.0; ) Alamofire/5.7.1").SetHeader("Accept-Encoding", "gzip").SetHeader("Connection", "close").SetHeader("Content-Type", "application/json")
}

func SklandRequest[T any](r *resty.Request, method, path string, vs ...any) (t T, _ error) {
	for i := 0; i < len(vs); i++ {
		switch v := vs[i].(type) {
		case AccountSkland:
			addSign(r, method, path, v)
		}
	}
	resp, respErr := r.SetError(&SKBaseResp[any]{}).SetResult(&SKBaseResp[T]{}).Execute(method, sklandAddr+path)
	if resp.StatusCode() == 405 {
		log.Println(string(resp.Body()))
		return t, fmt.Errorf("服务器被墙了！")
	}
	res, err := resty.ParseResp[*SKBaseResp[any], *SKBaseResp[T]](
		resp, respErr,
	)
	if err != nil {
		return t, fmt.Errorf("[skland] %w", err)
	}
	return res.Data, nil
}

func SklandRequestPlayerData(r *resty.Request, method, path string, vs ...any) (d string, _ error) {
	for i := 0; i < len(vs); i++ {
		switch v := vs[i].(type) {
		case AccountSkland:
			addSign(r, method, path, v)
		}
	}

	res, err := r.Execute(method, sklandAddr+path)
	if err != nil {
		return d, fmt.Errorf("[skland] %w", err)
	}
	if res.StatusCode() == 405 {
		log.Println(string(res.Body()))
		return d, fmt.Errorf("服务器被墙了！")
	}
	return string(res.Body()), nil
}
