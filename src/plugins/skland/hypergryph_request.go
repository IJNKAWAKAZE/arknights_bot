package skland

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/resty/v2"
)

var HypergryphAddr = "https://as.hypergryph.com"
var HypergryphAKAddr = "https://ak.hypergryph.com"

type HBaseResp[T any] struct {
	StatusCode *int   `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`

	Status *int   `json:"status"`
	Type   string `json:"type"`
	Msg    string `json:"msg"`

	Data T `json:"data,omitempty"`
}

type signHeaders struct {
	Platform  string `json:"platform"`
	Timestamp string `json:"timestamp"`
	DId       string `json:"dId"`
	VName     string `json:"vName"`
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

func HR() *resty.Request {
	return resty.R().SetHeader("User-Agent", viper.GetString("api.user_agent")).SetHeader("Accept-Encoding", "gzip")
}

func HypergryphRequest[T any](r *resty.Request, method, path string) (t T, _ error) {
	res, err := resty.ParseResp[*HBaseResp[any], *HBaseResp[T]](
		r.SetError(&HBaseResp[any]{}).SetResult(&HBaseResp[T]{}).Execute(method, HypergryphAddr+path),
	)
	if err != nil {
		return t, fmt.Errorf("[hypergryph] %w", err)
	}
	return res.Data, nil
}

func HypergryphASRequest(r *resty.Request, method, path string) (d string, _ error) {
	res, err := r.Execute(method, HypergryphAddr+path)
	if err != nil {
		return d, fmt.Errorf("[hypergryph] %w", err)
	}
	return string(res.Body()), nil
}

func HypergryphAKRequest(r *resty.Request, method, path string) (d string, _ error) {
	res, err := r.Execute(method, HypergryphAKAddr+path)
	if err != nil {
		return d, fmt.Errorf("[hypergryph] %w", err)
	}
	return string(res.Body()), nil
}
