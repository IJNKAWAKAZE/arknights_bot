package skland

import (
	"fmt"
	"github.com/starudream/go-lib/resty/v2"
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
	return resty.R().SetHeader("User-Agent", "Skland/1.5.1 (com.hypergryph.skland; build:100501001; Android 33; ) Okhttp/4.11.0").SetHeader("Accept-Encoding", "gzip")
}

func SklandRequest[T any](r *resty.Request, method, path string, vs ...any) (t T, _ error) {
	for i := 0; i < len(vs); i++ {
		switch v := vs[i].(type) {
		case AccountSkland:
			addSign(r, method, path, v)
		}
	}

	res, err := resty.ParseResp[*SKBaseResp[any], *SKBaseResp[T]](
		r.SetError(&SKBaseResp[any]{}).SetResult(&SKBaseResp[T]{}).Execute(method, sklandAddr+path),
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
	return string(res.Body()), nil
}
