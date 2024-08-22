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
	"sort"
	"strconv"
)

type DepotItem struct {
	Name   string `json:"name"`
	Count  string `json:"count"`
	Icon   string `json:"icon"`
	SortId int64  `json:"sortId"`
}

type ItemTable struct {
	Name   string `json:"name"`
	SortId int64  `json:"sortId"`
}

var itemMap = make(map[string]ItemTable)

func init() {
	resp, _ := http.Get(viper.GetString("api.item_table"))
	r, _ := io.ReadAll(resp.Body)
	gjson.ParseBytes(r).Get("items").ForEach(func(key, value gjson.Result) bool {
		itemMap[key.String()] = ItemTable{
			Name:   value.Get("name").String(),
			SortId: value.Get("sortId").Int(),
		}
		return true
	})
}

func Depot(r *gin.Engine) {
	r.GET("/depot", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Depot.tmpl")
		var depotItems []DepotItem
		var userAccount account.UserAccount
		var skAccount skland.Account
		userId, _ := strconv.ParseInt(c.Query("userId"), 10, 64)
		uid := c.Query("uid")
		sklandId := c.Query("sklandId")
		utils.GetAccountByUserIdAndSklandId(userId, sklandId).Scan(&userAccount)
		skAccount.Hypergryph.Token = userAccount.HypergryphToken
		skAccount.Skland.Token = userAccount.SklandToken
		skAccount.Skland.Cred = userAccount.SklandCred
		playerCultivate, err := skland.GetPlayerCultivate(uid, skAccount)
		if err != nil {
			log.Println(err)
			utils.WebC <- err
			return
		}
		for _, item := range playerCultivate.Items {
			if item.Count != "0" {
				var depotItem DepotItem
				depotItem.Name = itemMap[item.ID].Name
				depotItem.Count = item.Count
				count, _ := strconv.Atoi(depotItem.Count)
				if count >= 10000 {
					depotItem.Count = strconv.Itoa(count/10000) + "万"
				}
				depotItem.SortId = itemMap[item.ID].SortId
				// 图标
				paintingName := fmt.Sprintf("道具_带框_%s.png", depotItem.Name)
				m := utils.Md5(paintingName)
				path := "https://media.prts.wiki/thumb" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				depotItem.Icon = path + paintingName + "/75px-" + paintingName
				depotItems = append(depotItems, depotItem)
			}
		}
		sort.Slice(depotItems, func(i, j int) bool {
			return depotItems[i].SortId < depotItems[j].SortId
		})
		c.HTML(http.StatusOK, "Depot.tmpl", depotItems)
	})
}
