package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
)

// GetPlayerRedeem CDK兑换
func GetPlayerRedeem(token, cdk, channelId, uid string) (string, error) {
	if channelId == "1" {
		_, err := CheckToken(token)
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

	res, err := getPlayerRedeem(token, cdk, uid)
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

func getPlayerRedeem(token, cdk, uid string) (string, error) {

	u8Token, err := LoginHypergryph(token, uid)
	if err != nil {
		return "", fmt.Errorf("登录失败")
	}

	path := "/user/api/gift/exchange"
	headers := make(map[string]string)
	headers["Accept"] = "application/json"
	headers["X-Account-Token"] = token
	headers["X-Role-Token"] = u8Token
	req := HR().SetHeaders(headers).SetBody(gh.M{"giftCode": cdk})
	res, err := HypergryphAKRequest(req, "POST", path)
	if err != nil {
		return "", fmt.Errorf("发送兑换请求失败")
	}

	return res, nil
}
