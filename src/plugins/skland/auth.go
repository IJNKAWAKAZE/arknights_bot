package skland

import (
	"arknights_bot/config"
	"fmt"
	"github.com/starudream/go-lib/core/v2/gh"
	"github.com/tidwall/gjson"
	"log"
)

// Login 使用 token 登录
func Login(token, serverName string) (Account, error) {
	account := Account{}
	if token == "" {
		return account, fmt.Errorf("token is empty")
	}
	account.Hypergryph.Token = token
	isGlobal := serverName == serverGlobal

	appCode := "4ca99fa6b56cc2ba"
	if isGlobal {
		appCode = "6eb76d4e13aa36e6"
	}
	res, err := grantApp(token, appCode, isGlobal)
	if err != nil {
		return account, fmt.Errorf("grant app error: %w", err)
	}
	account.Hypergryph.Code = res.Code

	res1, err := authLoginByCode(res.Code, isGlobal)
	if err != nil {
		return account, fmt.Errorf("auth login by code error: %w", err)
	}
	u, err := checkToken(token, isGlobal)
	if err != nil {
		return account, fmt.Errorf("check token error: %w", err)
	}
	account.UserId = u.HgId
	account.Skland.Cred = res1.Cred
	account.Skland.Token = res1.Token
	return account, nil
}

// RefreshToken 刷新 token，cred 失效时重新登录并更新数据库
func RefreshToken(account Account, serverName string) (Account, error) {
	isGlobal := serverName == serverGlobal
	res, err := authRefresh(account.Skland.Cred, isGlobal)
	if err != nil {
		return account, fmt.Errorf("auth refresh error: %w", err)
	}
	account.Skland.Token = res.Token
	// 检查 cred 是否有效
	_, err = listPlayer(account.Skland, isGlobal)
	if err != nil {
		log.Println("cred失效，尝试重新登录。")
		if _, err = checkToken(account.Hypergryph.Token, isGlobal); err != nil {
			return account, err
		}
		account, err = Login(account.Hypergryph.Token, serverName)
		if err != nil {
			return account, err
		}
		// 更新 token
		config.DBEngine.Exec("update user_account set hypergryph_token = ?, skland_token = ?, skland_cred = ? where skland_id = ?",
			account.Hypergryph.Token, account.Skland.Token, account.Skland.Cred, account.UserId)
	}
	return account, nil
}

// ArknightsPlayers 获取明日方舟绑定角色
func ArknightsPlayers(skland AccountSkland, serverName string) ([]*Player, error) {
	playerList, err := listPlayer(skland, serverName == serverGlobal)
	if err != nil {
		return nil, err
	}
	for _, player := range playerList.List {
		if player.AppCode == "arknights" {
			return player.BindingList, nil
		}
	}
	return nil, nil
}

// loginHypergryph 官网登录获取 u8 token
func loginHypergryph(token, uid string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("token is empty")
	}

	reqGrantApp := HR().SetBody(gh.M{"type": 1, "token": token, "appCode": "be36d44aa36bfb5b"})
	grantAppToken, err := hgRawRequest(reqGrantApp, "POST", "/user/oauth2/v2/grant", hypergryphAddr)
	if err != nil {
		return "", fmt.Errorf("grant app error: %w", err)
	}

	reqU8Token := HR().SetBody(gh.M{"token": gjson.Parse(grantAppToken).Get("data.token").String(), "uid": uid})
	res, err := reqU8Token.Execute("POST", "https://binding-api-account-prod.hypergryph.com/account/binding/v1/u8_token_by_uid")
	if err != nil {
		return "", fmt.Errorf("get u8token error: %w", err)
	}
	u8Token := gjson.Parse(string(res.Body())).Get("data.token").String()

	reqLogin := HR().SetBody(gh.M{"share_by": "", "share_type": "", "source_from": "", "token": u8Token})
	resLogin, err := hgRawRequest(reqLogin, "POST", "/user/api/role/login", hypergryphAKAddr)
	if err != nil || resLogin == "" {
		return "", fmt.Errorf("登录失败：%w", err)
	}

	return u8Token, nil
}

// grantApp 获取 OAuth2 授权代码
func grantApp(token, appCode string, isGlobal bool) (*GrantAppData, error) {
	req := HR().SetBody(gh.M{"type": 0, "token": token, "appCode": appCode})
	return hgRequest[*GrantAppData](req, "POST", "/user/oauth2/v2/grant", isGlobal)
}

// authLoginByCode 获取 Cred
func authLoginByCode(code string, isGlobal bool) (*GenCredByCodeData, error) {
	req := SKR()
	if !isGlobal {
		req.SetHeader("did", did)
	}
	req.SetBody(gh.M{"kind": 1, "code": code})
	return skRequest[*GenCredByCodeData](req, "POST", "/web/v1/user/auth/generate_cred_by_code", isGlobal)
}

// authRefresh 刷新 auth
func authRefresh(cred string, isGlobal bool) (*AuthRefreshData, error) {
	req := SKR().SetHeader("cred", cred)
	return skRequest[*AuthRefreshData](req, "GET", "/api/v1/auth/refresh", isGlobal)
}

// listPlayer 获取绑定用户列表
func listPlayer(skland AccountSkland, isGlobal bool) (*ListPlayerData, error) {
	return skRequest[*ListPlayerData](SKR(), "GET", "/api/v1/game/player/binding", isGlobal, skland)
}

// checkToken 检查 token 有效性
func checkToken(token string, isGlobal bool) (*User, error) {
	req := HR().SetQueryParam("token", token)
	user, err := hgRequest[*User](req, "GET", "/user/info/v1/basic", isGlobal)
	if err != nil {
		return nil, fmt.Errorf("token已失效请重新登录！")
	}
	return user, err
}
