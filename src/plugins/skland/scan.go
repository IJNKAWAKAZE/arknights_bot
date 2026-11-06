package skland

import (
	"fmt"

	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
)

// 扫码登录状态码
const (
	ScanStatusConfirmed = 0   // 已确认，data.scanCode 有效
	ScanStatusWaiting   = 100 // 未扫码
	ScanStatusScanned   = 101 // 已扫码，待手机端确认
	ScanStatusExpired   = 102 // 二维码已过期
)

// ScanApp 可通过扫码登录的 App
type ScanApp struct {
	AppCode string `json:"appCode"`
	Name    string `json:"name"`
	IconUrl string `json:"iconUrl"`
}

// ScanLoginData 扫码登录二维码
type ScanLoginData struct {
	ScanId            string    `json:"scanId"`
	ScanUrl           string    `json:"scanUrl"`
	EnableScanAppList []ScanApp `json:"enableScanAppList"`
}

// GenScanLogin 生成扫码登录二维码
func GenScanLogin() (*ScanLoginData, error) {
	res, err := HR().Execute("POST", hypergryphAddr+"/general/v1/gen_scan/login")
	if err != nil {
		return nil, fmt.Errorf("生成二维码失败: %w", err)
	}
	r := gjson.ParseBytes(res.Body())
	if r.Get("status").Int() != 0 {
		return nil, fmt.Errorf("生成二维码失败: %s", r.Get("msg").String())
	}
	data := r.Get("data")
	if data.Get("scanId").String() == "" {
		return nil, fmt.Errorf("生成二维码失败: 未返回 scanId")
	}

	login := &ScanLoginData{
		ScanId:  data.Get("scanId").String(),
		ScanUrl: data.Get("scanUrl").String(),
	}
	for _, app := range data.Get("enableScanAppList").Array() {
		login.EnableScanAppList = append(login.EnableScanAppList, ScanApp{
			AppCode: app.Get("appCode").String(),
			Name:    app.Get("name").String(),
			IconUrl: app.Get("iconUrl").String(),
		})
	}
	return login, nil
}

// GetScanStatus 查询扫码状态，返回 (状态码, scanCode)
func GetScanStatus(scanId string) (int, string, error) {
	res, err := HR().
		SetQueryParam("scanId", scanId).
		Execute("GET", hypergryphAddr+"/general/v1/scan_status")
	if err != nil {
		return -1, "", fmt.Errorf("查询扫码状态失败: %w", err)
	}
	r := gjson.ParseBytes(res.Body())
	status := int(r.Get("status").Int())
	if status == ScanStatusConfirmed {
		return status, r.Get("data.scanCode").String(), nil
	}
	return status, "", nil
}

// TokenByScanCode 用扫码得到的 scanCode 换取鹰角账号 token
func TokenByScanCode(scanCode string) (string, error) {
	if scanCode == "" {
		return "", fmt.Errorf("scanCode 为空")
	}
	res, err := HR().
		SetBody(gh.M{"scanCode": scanCode}).
		Execute("POST", hypergryphAddr+"/user/auth/v1/token_by_scan_code")
	if err != nil {
		return "", fmt.Errorf("换取 token 失败: %w", err)
	}
	r := gjson.ParseBytes(res.Body())
	if r.Get("status").Int() != 0 {
		return "", fmt.Errorf("换取 token 失败: %s", r.Get("msg").String())
	}
	token := r.Get("data.token").String()
	if token == "" {
		return "", fmt.Errorf("换取 token 失败: 未返回 token")
	}
	return token, nil
}
