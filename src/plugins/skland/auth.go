package skland

import (
	bot "arknights_bot/config"
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
)

type GrantAppData struct {
	Uid  string `json:"uid"`
	Code string `json:"code"`
}

type GenCredByCodeData struct {
	UserId string `json:"userId"`
	Cred   string `json:"cred"`
	Token  string `json:"token"`
}

type AuthRefreshData struct {
	Token string `json:"token"`
}

type ListPlayerData struct {
	List []*PlayersByApp `json:"list"`
}

type PlayersByApp struct {
	AppCode     string    `json:"appCode"`
	AppName     string    `json:"appName"`
	DefaultUid  string    `json:"defaultUid"`
	BindingList []*Player `json:"bindingList"`
}

type Player struct {
	Uid             string `json:"uid"`
	ChannelName     string `json:"channelName"`
	ChannelMasterId string `json:"channelMasterId"`
	NickName        string `json:"nickName"`
	IsOfficial      bool   `json:"isOfficial"`
	IsDefault       bool   `json:"isDefault"`
	IsDelete        bool   `json:"isDelete"`
}

type GenTokenByUidData struct {
	Token string `json:"token"`
}

// Login 使用token登录
func Login(token string) (Account, error) {
	account := Account{}

	if token == "" {
		return account, fmt.Errorf("token is empty")
	}
	account.Hypergryph.Token = token

	res, err := grantApp(token, "4ca99fa6b56cc2ba")
	if err != nil {
		return account, fmt.Errorf("grant app error: %w", err)
	}
	account.Hypergryph.Code = res.Code

	res1, err := authLoginByCode(res.Code)
	if err != nil {
		return account, fmt.Errorf("auth login by code error: %w", err)
	}
	account.UserId = res1.UserId
	account.Skland.Cred = res1.Cred
	account.Skland.Token = res1.Token
	return account, nil
}

// 获取 OAuth2 授权代码
func grantApp(token string, code string) (*GrantAppData, error) {
	req := HR().SetBody(gh.M{"type": 0, "token": token, "appCode": code})
	return HypergryphRequest[*GrantAppData](req, "POST", "/user/oauth2/v2/grant")
}

// 获取Cred
func authLoginByCode(code string) (*GenCredByCodeData, error) {
	req := SKR().SetBody(gh.M{"kind": 1, "code": code})
	return SklandRequest[*GenCredByCodeData](req, "POST", "/api/v1/user/auth/generate_cred_by_code")
}

// RefreshToken 刷新 token
func RefreshToken(account Account) (Account, error) {
	b, err := CheckUser(account.Skland.Cred)
	if err != nil {
		return account, err
	}
	if b {
		res, err := authRefresh(account.Skland.Cred)
		if err != nil {
			return account, fmt.Errorf("auth refresh error: %w", err)
		}
		account.Skland.Token = res.Token
	} else {
		err = CheckToken(account.Hypergryph.Token)
		if err != nil {
			return account, err
		}
		account, err = Login(account.Hypergryph.Token)
		if err != nil {
			return account, err
		}
		// 更新token
		bot.DBEngine.Exec("update user_account set hypergryph_token = ?, skland_token = ?, skland_cred = ? where skland_id = ?", account.Hypergryph.Token, account.Skland.Token, account.Skland.Cred, account.UserId)
	}
	return account, nil
}

// CheckUser 检查用户cred有效性
func CheckUser(cred string) (bool, error) {
	res, err := SKR().SetHeader("cred", cred).Get(sklandAddr + "/api/v1/user/check")
	if err != nil {
		return false, err
	}
	if res.StatusCode() == 401 {
		return false, fmt.Errorf("[skland] %s", "用户未登录！")
	}
	return true, nil
}

// CheckToken 检查token有效性
func CheckToken(token string) error {
	req := SKR().SetQueryParam("token", token)
	_, err := HypergryphRequest[any](req, "GET", "/user/info/v1/basic")
	if err != nil {
		return fmt.Errorf("token已失效请重新登录！")
	}
	return err
}

// CheckBToken 检查BToken有效性
func CheckBToken(token string) error {
	req := HR().SetBody(gh.MS{"token": token})
	_, err := HypergryphRequest[any](req, "POST", "/u8/user/info/v1/basic")
	if err != nil {
		return fmt.Errorf("btoken已失效请重新登录！")
	}
	return err
}

// 刷新 auth
func authRefresh(cred string) (*AuthRefreshData, error) {
	req := SKR().SetHeader("cred", cred)
	return SklandRequest[*AuthRefreshData](req, "GET", "/api/v1/auth/refresh")
}

// 获取绑定用户列表
func listPlayer(skland AccountSkland) (*ListPlayerData, error) {
	return SklandRequest[*ListPlayerData](SKR(), "GET", "/api/v1/game/player/binding", skland)
}

// ArknightsPlayers 获取明日方舟绑定角色
func ArknightsPlayers(skland AccountSkland) ([]*Player, error) {
	var players []*Player
	playerList, err := listPlayer(skland)
	if err != nil {
		return players, err
	}
	for _, player := range playerList.List {
		if player.AppCode == "arknights" {
			players = player.BindingList
		}
	}
	return players, nil
}
