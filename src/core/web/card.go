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
	"net/url"
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

var ignoreChar = make(map[string]string)

func init() {
	ignoreChar["npc_001_doctor"] = "npc_001_doctor"
	ignoreChar["npc_003_kalts"] = "npc_003_kalts"
	ignoreChar["npc_010_chen"] = "npc_010_chen"
	ignoreChar["npc_2005_wywu"] = "npc_2005_wywu"
	ignoreChar["npc_2006_fmzuki"] = "npc_2006_fmzuki"
	ignoreChar["char_513_apionr"] = "char_513_apionr"
	ignoreChar["char_511_asnipe"] = "char_511_asnipe"
	ignoreChar["char_510_amedic"] = "char_510_amedic"
	ignoreChar["char_508_aguard"] = "char_508_aguard"
	ignoreChar["char_509_acast"] = "char_509_acast"
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
	//playerCard.Resume = userPlayer.Resume
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
	avatarId := playerData.Status.Avatar.Id
	if strings.HasPrefix(avatarId, "char") {
		playerCard.Avatar = fmt.Sprintf("https://web.hycdn.cn/arknights/game/assets/char_skin/avatar/%s.png", url.QueryEscape(avatarId))
	} else {
		// 头像
		paintingName := fmt.Sprintf("%s.png", strings.ToUpper(avatarId[:1])+avatarId[1:])
		m := utils.Md5(paintingName)
		path := "https://media.prts.wiki/thumb" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
		playerCard.Avatar = path + paintingName + "/80px-" + paintingName
	}
	playerCard.Resume = playerData.Status.Resume
	playerCard.CharCnt = len(playerData.Chars)
	playerCard.NationList = getNationList(playerData)
	if _, has := charMap["char_1001_amiya2"]; has {
		playerCard.CharCnt -= 1
	}
	if _, has := charMap["char_1037_amiya3"]; has {
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
	if strings.Contains(skinId, "amiya3") {
		secretaryName = "阿米娅(医疗)"
	}
	operator := utils.GetOperatorByName(secretaryName)
	enName := operator.NameEn
	if !strings.Contains(skinId, "@") {
		if strings.Contains(skinId, "#1") {
			skinUrl = operator.Skins[0].Url
		} else if strings.Contains(skinId, "#2") && (secretaryName == "阿米娅(近卫)" || secretaryName == "阿米娅(医疗)") {
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

	var m = make(map[string][]string)
	resp, _ := http.Get(viper.GetString("api.nation_table"))
	r, _ := io.ReadAll(resp.Body)
	gjson.ParseBytes(r).Get("groupList").ForEach(func(key, value gjson.Result) bool {
		charList := value.Get("charList").Array()
		for _, c := range charList {
			if _, has := ignoreChar[c.Get("charId").String()]; !has {
				m[key.String()] = append(m[key.String()], c.Get("charId").String())
			}
		}
		return true
	})
	defer resp.Body.Close()
	for i, nation := range nationList {
		count := 0
		key := nation.Name
		for _, char := range m[key] {
			if _, has := playerData.CharInfoMap[char]; has {
				count++
			}
		}
		if count == 0 {
			nationList[i].Flag = -1
		} else if len(m[key]) > count {
			nationList[i].Flag = 0
		} else if len(m[key]) == count {
			nationList[i].Flag = 1
		}

	}
	return nationList
}
