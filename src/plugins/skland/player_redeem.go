package skland

import (
	"fmt"
	"github.com/eddieivan01/nic"
	"github.com/tidwall/gjson"
	"io"
	"log"
)

// GetPlayerRedeem CDK兑换
func GetPlayerRedeem(token, cdk string) (string, error) {
	err := checkToken(token)
	if err != nil {
		log.Println(err)
		return err.Error(), err
	}

	res, err := getPlayerRedeem(token, cdk)
	if err != nil {
		return "", err
	}
	code := gjson.Get(res, "code").String()
	msg := gjson.Get(res, "msg").String()
	if code != "" && code != "0" {
		return msg, nil
	}
	return "", nil
}

func getPlayerRedeem(token, cdk string) (string, error) {
	session := nic.NewSession()
	resp, err := session.Get(HypergryphAKAddr+"/user/api/gift/getExchangeLog?token="+token+"&channelId=1", nil)
	if err != nil {
		return "", fmt.Errorf("获取csrf_token失败")
	}
	post, err := session.Post(HypergryphAKAddr+"/user/api/gift/exchange", &nic.H{
		JSON: nic.KV{
			"token":     token,
			"giftCode":  cdk,
			"channelId": 1,
		},
		Headers: nic.KV{
			"authority":          "ak.hypergryph.com",
			"Accept":             "application/json",
			"Accept-Encoding":    "gzip, deflate, br",
			"Accept-Language":    "zh-CN,zh;q=0.8",
			"Content-Length":     "80",
			"Content-Type":       "application/json;charset=UTF-8",
			"Origin":             "https://ak.hypergryph.com",
			"Referer":            "https://ak.hypergryph.com/user/exchangeGift",
			"Sec-Ch-Ua":          "\"Not A(Brand\";v=\"99\", \"Brave\";v=\"121\", \"Chromium\";v=\"121\"",
			"Sec-Ch-Ua-Mobile":   "?0",
			"Sec-Ch-Ua-Platform": "\"Windows\"",
			"Sec-Fetch-Dest":     "empty",
			"Sec-Fetch-Mode":     "cors",
			"Sec-Fetch-Site":     "same-origin",
			"Sec-Gpc":            "1",
			"User-Agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36",
			"X-Csrf-Token":       resp.Cookies()[1].Value,
		},
	})
	if err != nil {
		return "", fmt.Errorf("发送兑换请求失败")
	}
	read, _ := io.ReadAll(post.Body)
	return string(read), nil
}
