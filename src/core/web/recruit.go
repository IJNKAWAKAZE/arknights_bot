package web

import (
	"arknights_bot/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"sort"
	"strings"
)

type RecruitList struct {
	Tags      []string         `json:"tags"`
	Operators []utils.Operator `json:"operators"`
}

func Recruit(r *gin.Engine) {
	r.GET("/recruit", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Recruit.tmpl")
		tags := strings.Split(c.Query("tags"), " ")
		var recruitList []RecruitList
		recruitOperatorList := utils.GetRecruitOperatorList()
		var tagList [][]string
		n := len(tags)
		nBit := 1 << n
		for i := 1; i < nBit; i++ {
			var ts []string
			for j := 0; j < n; j++ {
				tmp := 1 << j
				if (tmp & i) != 0 {
					ts = append(ts, tags[j])
				}
			}
			tagList = append(tagList, ts)
		}
		sort.Slice(tagList, func(i, j int) bool {
			if strings.Contains(strings.Join(tagList[i], " "), "高资") {
				return true
			}
			return len(tagList[i]) > len(tagList[j])
		})
		for _, t := range tagList {
			var recruit RecruitList
			tl := strings.Join(t, " ")
			replacer := strings.NewReplacer("机械", "支援机械", "资深", "资深干员", "高资", "高级资深干员")
			recruit.Tags = strings.Split(replacer.Replace(tl), " ")
			if strings.Contains(tl, "高资") && strings.Contains(tl, "资深") {
				continue
			}
			for _, operator := range recruitOperatorList {
				opTags := operator.Tags + " " + operator.ProfessionZH + " " + operator.Position
				if operator.Rarity == 5 {
					opTags += " 高资"
				}
				if operator.Rarity == 4 {
					opTags += " 资深"
				}
				b := true
				for _, tag := range t {
					if !strings.Contains(opTags, tag) {
						b = false
					}
					if operator.Rarity == 5 {
						if !strings.Contains(tl, "高资") {
							b = false
						}
					}
				}
				if b {
					recruit.Operators = append(recruit.Operators, operator)
				}
			}
			if len(recruit.Operators) > 0 {
				sort.Slice(recruit.Operators, func(i, j int) bool {
					return recruit.Operators[i].Rarity > recruit.Operators[j].Rarity
				})
				recruitList = append(recruitList, recruit)
			}
		}
		c.HTML(http.StatusOK, "Recruit.tmpl", recruitList)
	})
}
