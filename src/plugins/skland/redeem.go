package skland

import (
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
)

// GetPlayerRedeem CDK兑换
func GetPlayerRedeem(token, cdk, uid string) (string, error) {
	if _, err := checkToken(token, false); err != nil {
		log.Println(err)
		return err.Error(), err
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

// getPlayerRedeem 发送兑换请求
func getPlayerRedeem(token, cdk, uid string) (string, error) {
	u8Token, err := loginHypergryph(token, uid)
	if err != nil {
		return "", fmt.Errorf("登录失败")
	}
	req := HR().SetHeaders(map[string]string{
		"Accept":          "application/json",
		"X-Account-Token": token,
		"X-Role-Token":    u8Token,
	}).SetBody(gh.M{"giftCode": cdk})
	res, err := hgRawRequest(req, "POST", "/user/api/gift/exchange", hypergryphAKAddr)
	if err != nil {
		return "", fmt.Errorf("发送兑换请求失败: %w", err)
	}
	return res, nil
}
