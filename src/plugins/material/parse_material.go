package material

import (
	"arknights_bot/utils"
	"fmt"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

type Material struct {
	Name   string  `json:"name"`
	Stages []Stage `json:"stages"`
}

type Stage struct {
	ZoneName          string `json:"zoneName"`          // 关卡
	Code              string `json:"code"`              // 编码
	Name              string `json:"name"`              // 主材料名称
	Icon              string `json:"icon"`              // 主材料图标
	KnockRating       string `json:"knockRating"`       // 主产物掉率
	ApExpect          string `json:"apExpect"`          // 期望理智
	SecondaryItem     string `json:"SecondaryItem"`     // 副产物名称
	SecondaryItemIcon string `json:"SecondaryItemIcon"` // 副产物图标
	StageEfficiency   string `json:"stageEfficiency"`   // 关卡效率
}

var itemMap = make(map[string]string)

func ParseMaterial(name string) Material {
	var material Material
	var stages []Stage
	if len(itemMap) == 0 {
		res, _ := http.Get(viper.GetString("api.item"))
		read, _ := io.ReadAll(res.Body)
		defer res.Body.Close()
		j := gjson.ParseBytes(read)
		for _, item := range j.Get("data").Array() {
			itemMap[item.Get("itemId").String()] = item.Get("itemName").String()
		}
	}
	res, err := http.Get(viper.GetString("api.stage_result"))
	if err != nil {
		return material
	}
	read, err := io.ReadAll(res.Body)
	if err != nil {
		return material
	}
	defer res.Body.Close()
	j := gjson.ParseBytes(read)
	for _, d := range j.Get("data.recommendedStageList").Array() {
		itemType := d.Get("itemType").String()
		if itemType == name {
			material.Name = itemType
			for _, item := range d.Get("stageResultList").Array() {
				var stage Stage
				stage.ZoneName = item.Get("zoneName").String()
				stage.Code = item.Get("stageCode").String()
				stage.Name = item.Get("itemName").String()
				paintingName := fmt.Sprintf("道具_%s.png", stage.Name)
				m := utils.Md5(paintingName)
				path := "https://prts.wiki" + fmt.Sprintf("/images/thumb/%s/%s/", m[:1], m[:2])
				pic := path + paintingName + "/75px-" + paintingName
				stage.Icon = pic
				stage.ApExpect = fmt.Sprintf("%.1f", item.Get("apExpect").Float())
				stage.KnockRating = fmt.Sprintf("%.1f%%", item.Get("knockRating").Float()*100)
				stage.SecondaryItem = itemMap[item.Get("secondaryItemId").String()]
				if stage.SecondaryItem != "" {
					paintingName := fmt.Sprintf("道具_%s.png", stage.SecondaryItem)
					m := utils.Md5(paintingName)
					path := "https://prts.wiki" + fmt.Sprintf("/images/thumb/%s/%s/", m[:1], m[:2])
					pic := path + paintingName + "/75px-" + paintingName
					stage.SecondaryItemIcon = pic
				}
				stage.StageEfficiency = fmt.Sprintf("%.1f%%", item.Get("stageEfficiency").Float()*100)
				stages = append(stages, stage)
			}
			break
		}
	}
	material.Stages = stages
	return material
}
