package skland

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"github.com/starudream/go-lib/core/v2/codec/json"
	"github.com/starudream/go-lib/core/v2/utils/structutil"
	"github.com/starudream/go-lib/resty/v2"
	"net/url"
	"strconv"
	"time"
)

func addSign(r *resty.Request, method, path string, skland AccountSkland) {
	ts := strconv.FormatInt(time.Now().Unix()-7, 10)

	headers := signHeaders{Platform: "1", Timestamp: ts, DId: "743a446c83032899", VName: "1.21.0"}

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
