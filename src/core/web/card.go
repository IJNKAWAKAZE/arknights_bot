package web

import (
	"arknights_bot/plugins/account"
	"arknights_bot/plugins/skland"
	"arknights_bot/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type PlayerCard struct {
	Name              string   `json:"name"`
	Uid               string   `json:"uid"`
	ServerName        string   `json:"serverName"`
	Resume            string   `json:"resume"`
	Level             int      `json:"level"`
	RegTime           int      `json:"regTime"`
	MainStageProgress string   `json:"mainStageProgress"`
	Avatar            string   `json:"avatar"`
	SecretaryName     string   `json:"secretaryName"`
	SecretaryEnName   string   `json:"secretaryEnName"`
	Secretary         string   `json:"secretary"`
	CharCnt           int      `json:"charCnt"`
	FurnitureCnt      int      `json:"furnitureCnt"`
	SkinCnt           int      `json:"skinCnt"`
	NationList        []Nation `json:"nationList"`
	AssistChars       []struct {
		Name            string `json:"name"`
		CharID          string `json:"charId"`
		SkinID          string `json:"skinId"`
		Level           int    `json:"level"`
		EvolvePhase     int    `json:"evolvePhase"`
		PotentialRank   int    `json:"potentialRank"`
		SkillID         string `json:"skillId"`
		MainSkillLvl    int    `json:"mainSkillLvl"`
		SpecializeLevel int    `json:"specializeLevel"`
		Equip           struct {
			ID           string `json:"id"`
			Level        int    `json:"level"`
			Name         string `json:"name"`
			TypeIcon     string `json:"typeIcon"`
			ShiningColor string `json:"shiningColor"`
		} `json:"equip"`
	} `json:"assistChars"`
}

type Nation struct {
	Name string `json:"name"`
	Flag int64  `json:"flag"`
}

func Card(r *gin.Engine) {
	r.GET("/card", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Card.tmpl")
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		playerCard, err := cardData(userId, sklandId, uid)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "Card.tmpl", playerCard)
	})

	r.GET("/oldCard", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/OldCard.tmpl")
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		playerCard, err := cardData(userId, sklandId, uid)
		if err != nil {
			return
		}
		c.HTML(http.StatusOK, "OldCard.tmpl", playerCard)
	})
}

func cardData(userId int64, sklandId, uid string) (PlayerCard, error) {
	var playerCard PlayerCard
	var userAccount account.UserAccount
	var skAccount skland.Account
	utils.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
	var userPlayer account.UserPlayer
	utils.GetPlayerByUserId(userAccount.UserNumber, uid).Scan(&userPlayer)
	playerCard.ServerName = userPlayer.ServerName
	playerCard.Resume = userPlayer.Resume
	skAccount.Hypergryph.Token = userAccount.HypergryphToken
	skAccount.Skland.Token = userAccount.SklandToken
	skAccount.Skland.Cred = userAccount.SklandCred
	playerData, skAccount, err := skland.GetPlayerInfo(uid, skAccount)
	if err != nil {
		log.Println(err)
		return playerCard, fmt.Errorf("获取名片信息失败")
	}

	charMap := playerData.CharInfoMap
	secretaryName := charMap[playerData.Status.Secretary.CharID].Name
	skinUrl, enName, err := getSkinUrl(secretaryName, playerData.Status.Secretary.SkinID)
	if err != nil {
		return playerCard, err
	}
	playerCard.Secretary = skinUrl
	playerCard.SecretaryName = secretaryName
	playerCard.SecretaryEnName = enName
	playerCard.Name = playerData.Status.Name
	playerCard.Uid = playerData.Status.UID
	playerCard.Level = playerData.Status.Level
	playerCard.RegTime = playerData.Status.RegisterTs
	playerCard.MainStageProgress = playerData.StageInfoMap[playerData.Status.MainStageProgress].Code
	playerCard.Avatar = playerData.Status.Secretary.SkinID
	playerCard.CharCnt = len(playerData.Chars)
	playerCard.NationList = getNationList(playerData)
	if _, has := charMap["char_1001_amiya2"]; has {
		playerCard.CharCnt -= 1
	}
	playerCard.SkinCnt = len(playerData.Skins)
	playerCard.FurnitureCnt = playerData.Building.Furniture.Total
	playerCard.AssistChars = playerData.AssistChars
	for i, char := range playerCard.AssistChars {
		name := playerData.CharInfoMap[char.CharID].Name
		equipId := playerCard.AssistChars[i].Equip.ID
		playerCard.AssistChars[i].Name = name
		playerCard.AssistChars[i].Equip.Name = strings.ToUpper(playerData.EquipmentInfoMap[equipId].TypeIcon)
		playerCard.AssistChars[i].Equip.TypeIcon = playerData.EquipmentInfoMap[equipId].TypeIcon
		playerCard.AssistChars[i].Equip.ShiningColor = playerData.EquipmentInfoMap[equipId].ShiningColor
		playerCard.AssistChars[i].MainSkillLvl = playerCard.AssistChars[i].MainSkillLvl + playerCard.AssistChars[i].SpecializeLevel
	}
	return playerCard, nil
}

func getSkinUrl(secretaryName, skinId string) (string, string, error) {
	skinUrl := ""
	if strings.Contains(skinId, "amiya2") {
		secretaryName = "阿米娅(近卫)"
	}
	operator := utils.GetOperatorByName(secretaryName)
	enName := operator.NameEn
	if !strings.Contains(skinId, "@") {
		if strings.Contains(skinId, "#1") {
			skinUrl = operator.Skins[0].Url
		} else if strings.Contains(skinId, "#2") && secretaryName == "阿米娅(近卫)" {
			skinUrl = operator.Skins[0].Url
		} else if strings.Contains(skinId, "#2") {
			skinUrl = operator.Skins[1].Url
		} else if strings.Contains(skinId, "#1+") {
			skinUrl = operator.Skins[2].Url
		}
		return skinUrl, enName, nil
	}

	resp, err := http.Get(viper.GetString("api.skin_table"))
	if err != nil {
		log.Println(err)
		return skinUrl, enName, err
	}
	dataByte, _ := io.ReadAll(resp.Body)
	skinTable := gjson.ParseBytes(dataByte)
	defer resp.Body.Close()

	for _, skin := range skinTable.Get("charSkins").Array() {
		skin.ForEach(func(key, value gjson.Result) bool {
			if key.String() == skinId {
				skinName := value.Get("displaySkin.skinName").String()
				for _, sk := range operator.Skins {
					if sk.Name == skinName {
						skinUrl = sk.Url
					}
				}
				return false
			}
			return true
		})
	}
	return skinUrl, enName, nil
}

func getNationList(playerData *skland.PlayerData) []Nation {
	nationList := []Nation{
		{Name: "rhodes"},
		{Name: "lungmen"},
		{Name: "yan"},
		{Name: "egir"},
		{Name: "bolivar"},
		{Name: "columbia"},
		{Name: "higashi"},
		{Name: "iberia"},
		{Name: "kazimierz"},
		{Name: "kjerag"},
		{Name: "laterano"},
		{Name: "leithanien"},
		{Name: "minos"},
		{Name: "rim"},
		{Name: "sami"},
		{Name: "sargon"},
		{Name: "siracusa"},
		{Name: "ursus"},
		{Name: "victoria"},
		{Name: "kazdel"},
		{Name: "followers"}}

	var m = make(map[string]int)
	var m1 = make(map[string]int)
	resp, _ := http.Get("https://raw.githubusercontent.com/Kengxxiao/ArknightsGameData/master/zh_CN/gamedata/art/handbookpos_table.json")
	r, _ := io.ReadAll(resp.Body)
	gjson.ParseBytes(r).Get("groupList").ForEach(func(key, value gjson.Result) bool {
		m[key.String()] = len(value.Get("charList").Array())
		return true
	})
	for _, v := range playerData.CharInfoMap {
		key := v.NationID
		if v.NationID == "" {
			key = "followers"
		}
		m1[key] += 1
	}
	defer resp.Body.Close()
	for i, nation := range nationList {
		key := nation.Name
		if m1[nation.Name] == 0 {
			nationList[i].Flag = -1
		} else if m[key] > m1[key] {
			nationList[i].Flag = 0
		} else if m[key] == m1[key] {
			nationList[i].Flag = 1
		}
	}
	return nationList
}
