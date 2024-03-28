package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
)

// GetPlayerRedeem CDK兑换
func GetPlayerRedeem(token, cdk, channelId string) (string, error) {
	if channelId == "1" {
		err := CheckToken(token)
		if err != nil {
			log.Println(err)
			return err.Error(), err
		}
	} else if channelId == "2" {
		err := CheckBToken(token)
		if err != nil {
			log.Println(err)
			return err.Error(), err
		}
	}

	res, err := getPlayerRedeem(token, cdk, channelId)
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

func getPlayerRedeem(token, cdk, channelId string) (string, error) {
	path := "/user/api/gift/exchange"
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	resp, err := http.Get(HypergryphAKAddr + path)
	if err != nil {
		return "", fmt.Errorf("获取csrf_token失败")
	}
	csrfToken := ""
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "csrf_token" {
			csrfToken = cookie.Value
		}
	}
	headers["X-Csrf-Token"] = csrfToken
	req1 := SKR().SetHeaders(headers).SetBody(gh.M{"token": token, "giftCode": cdk, "channelId": channelId}).SetCookies(resp.Cookies())
	res, err := HypergryphAKRequest(req1, "POST", path)
	if err != nil {
		return "", fmt.Errorf("发送兑换请求失败")
	}

	return res, nil
}
