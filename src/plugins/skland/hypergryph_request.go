package skland

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/starudream/go-lib/core/v2/codec/json"
	"github.com/starudream/go-lib/core/v2/utils/structutil"
	"github.com/starudream/go-lib/resty/v2"
	"net/url"
	"strconv"
	"time"
)

type HBaseResp[T any] struct {
	StatusCode *int   `json:"statusCode"`
	Error      string `json:"error"`
	Message    string `json:"message"`

	Status *int   `json:"status"`
	Type   string `json:"type"`
	Msg    string `json:"msg"`

	Data T `json:"data,omitempty"`
}

type SKBaseResp[T any] struct {
	Code    *int   `json:"code"`
	Message string `json:"message"`
	Data    T      `json:"data,omitempty"`
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

func (t *SKBaseResp[T]) IsSuccess() bool {
	return t != nil && t.Code != nil && *t.Code == 0
}

func (t *HBaseResp[T]) String() string {
	if t != nil && t.StatusCode != nil {
		return fmt.Sprintf("status: %d, error: %s, message: %s", *t.StatusCode, t.Error, t.Message)
	} else if t != nil && t.Status != nil {
		return fmt.Sprintf("status: %d, type: %s, msg: %s", *t.Status, t.Type, t.Msg)
	}
	return "<nil>"
}

func (t *SKBaseResp[T]) String() string {
	return fmt.Sprintf("code: %d, message: %s", *t.Code, t.Message)
}

func HR() *resty.Request {
	return resty.R().SetHeader("User-Agent", HypergryphUserAgent).SetHeader("Accept-Encoding", "gzip")
}

func SKR() *resty.Request {
	return resty.R().SetHeader("User-Agent", SklandUserAgent).SetHeader("Accept-Encoding", "gzip")
}

func IsUnauthorized(err error) bool {
	re, ok := resty.AsRespErr(err)
	if ok {
		return re.StatusCode() == 401
	}
	return false
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

func SklandRequest[T any](r *resty.Request, method, path string, vs ...any) (t T, _ error) {
	for i := 0; i < len(vs); i++ {
		switch v := vs[i].(type) {
		case AccountSkland:
			addSign(r, method, path, v)
		}
	}

	res, err := resty.ParseResp[*SKBaseResp[any], *SKBaseResp[T]](
		r.SetError(&SKBaseResp[any]{}).SetResult(&SKBaseResp[T]{}).Execute(method, SklandAddr+path),
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

	res, err := r.Execute(method, SklandAddr+path)
	if err != nil {
		return d, fmt.Errorf("[skland] %w", err)
	}
	return string(res.Body()), nil
}

func addSign(r *resty.Request, method, path string, skland AccountSkland) {
	ts := strconv.FormatInt(time.Now().Unix(), 10)

	headers := signHeaders{Platform: Platform, Timestamp: ts, DId: DId, VName: VName}

	r.SetHeaders(tom(headers))

	_, signature := sign(headers, method, path, skland.Token, r.QueryParam, r.Body)

	r.SetHeader("cred", skland.Cred)
	r.SetHeader("sign", signature)
}

func tom(s any) map[string]string {
	t := structutil.New(s)
	t.TagName = "json"
	m := map[string]string{}
	for k, v := range t.Map() {
		m[k] = v.(string)
	}
	return m
}

func sign(headers signHeaders, method, path, token string, query url.Values, body any) (string, string) {
	str := query.Encode()
	if method != "GET" {
		str = json.MustMarshalString(body)
	}

	content := path + str + headers.Timestamp + json.MustMarshalString(headers)

	b1 := hmac256(token, content)
	s1 := hex.EncodeToString(b1)
	b2 := md5.Sum([]byte(s1))
	s2 := hex.EncodeToString(b2[:])

	return content, s2
}

func hmac256(key, content string) []byte {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(content))
	return h.Sum(nil)
}
