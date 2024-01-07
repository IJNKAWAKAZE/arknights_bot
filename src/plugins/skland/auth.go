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

type User struct {
	User *UserInfo `json:"user"`
}

type UserInfo struct {
	Id       string `json:"id"`
	Nickname string `json:"nickname"`
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

// Login 使用token登录
func Login(token string) (Account, error) {
	account := Account{}

	if token == "" {
		return account, fmt.Errorf("token is empty")
	}
	account.Hypergryph.Token = token

	res, err := grantApp(token, AppCodeSKLAND)
	if err != nil {
		return account, fmt.Errorf("grant app error: %w", err)
	}
	account.Hypergryph.Code = res.Code

	res1, err := authLoginByCode(res.Code)
	if err != nil {
		return account, fmt.Errorf("auth login by code error: %w", err)
	}
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
func RefreshToken(uid string, account Account) (Account, error) {
	_, err := getUser(account.Skland)
	if err == nil {
		return account, nil
	}
	if !IsUnauthorized(err) {
		return account, fmt.Errorf("get user error: %w", err)
	}

	res, err := authRefresh(account.Skland.Cred)
	if err != nil {
		return account, fmt.Errorf("auth refresh error: %w", err)
	}
	account.Skland.Token = res.Token

	_, err = getUser(account.Skland)
	if err != nil {
		if !IsUnauthorized(err) {
			return account, fmt.Errorf("get user error: %w", err)
		}
		account, err = Login(account.Hypergryph.Token)
		if err != nil {
			return account, err
		}
	}
	// 查询更新用户
	var userNumber string
	result := bot.DBEngine.Raw("select user_number from user_player where uid = ?", uid).Scan(&userNumber)

	if result.RowsAffected > 0 {
		// 更新token
		bot.DBEngine.Exec("update user_account set hypergryph_token = ?, skland_token = ?, skland_cred = ? where user_number = ?", account.Hypergryph.Token, account.Skland.Token, account.Skland.Cred, userNumber)
	}
	return account, nil
}

// 获取用户信息
func getUser(skland AccountSkland) (*User, error) {
	return SklandRequest[*User](SKR(), "GET", "/api/v1/user", skland)
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

// ArknihghtsPlayers 获取明日方舟绑定角色
func ArknihghtsPlayers(skland AccountSkland) ([]*Player, error) {
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
