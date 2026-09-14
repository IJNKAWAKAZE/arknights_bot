package operator

import (
	"arknights_bot/utils/hashutil"
	"arknights_bot/utils/model"
	"arknights_bot/utils/search"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"html/template"
	"net/http"
	"strconv"
	"strings"
)

type Operator struct {
	OP               model.Operator   `json:"op"`               // 基本信息
	Painting         string           `json:"painting"`         // 立绘
	AttackRange      string           `json:"attackRange"`      // 攻击范围
	ProfessionBranch ProfessionBranch `json:"professionBranch"` // 职业分支
	Potentials       []Potential      `json:"potentials"`       // 潜能
	Talents          []Talent         `json:"talents"`          // 天赋
	BuildingSkills   []BuildingSkill  `json:"buildingSkills"`   // 基建技能
	Skills           []Skill          `json:"skills"`           // 技能
}

type ProfessionBranch struct {
	Name string `json:"name"` // 名称
	Pic  string `json:"pic"`  // 图片
	Desc string `json:"desc"` // 描述
}

type Potential struct {
	Rank int    `json:"rank"` // 等级
	Desc string `json:"desc"` // 描述
}

type Talent struct {
	Evolve string        `json:"evolve"` // 精英级别
	Name   template.HTML `json:"name"`   // 名称
	Desc   template.HTML `json:"desc"`   // 描述
}

type BuildingSkill struct {
	Evolve string `json:"evolve"` // 精英级别
	Icon   string `json:"icon"`   // 图标
	Name   string `json:"name"`   // 名称
	Desc   string `json:"desc"`   // 描述
}

type Skill struct {
	Icon       string          `json:"icon"`       // 图标
	Name       string          `json:"name"`       // 名称
	Desc       template.HTML   `json:"desc"`       // 描述
	SkillRange template.HTML   `json:"skillRange"` // 技能范围
	SpType     []template.HTML `json:"spType"`     // 回费类型
	SpInit     string          `json:"spInit"`     // 初始费用
	SpCost     string          `json:"spCost"`     // 所需费用
	Duration   string          `json:"duration"`   // 持续时间
}

// ParseOperator 解析干员数据
func ParseOperator(name string) Operator {
	var operator Operator
	api := viper.GetString("api.wiki")
	response, _ := http.Get(api + name)
	op := search.GetOperatorByName(name)
	if op.Name != "" {
		operator.OP = op
		operator.Painting = op.Skins[0].Url
		if op.Rarity > 2 && len(op.Skins) > 1 {
			operator.Painting = op.Skins[1].Url
		}
		if op.Name == "阿米娅(近卫)" || op.Name == "阿米娅(医疗)" {
			operator.Painting = op.Skins[0].Url
		}
		doc, _ := goquery.NewDocumentFromReader(response.Body)

		// 职业分支
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "特性" {
				selection.Parent().NextFilteredUntil(".wikitable", ".mw-heading").Each(func(j int, selection *goquery.Selection) {
					// 剔除 wiki 的“术语: xxx”悬浮说明（隐藏节点），避免特性描述里混入大段注释
					selection.Find(".mc-tooltips span[data-append-to]").Remove()
					tds := selection.Find("td")
					operator.ProfessionBranch.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
					paintingName := fmt.Sprintf("职业分支图标_%s.png", operator.ProfessionBranch.Name)
					m := hashutil.Md5(paintingName)
					path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
					operator.ProfessionBranch.Pic = path + paintingName
					operator.ProfessionBranch.Desc = strings.ReplaceAll(tds.Eq(1).Text(), "\n", "")
				})
			}
		})

		// 信赖加成
		doc.Find(".char-base-attr-table tr").Each(func(i int, selection *goquery.Selection) {
			selection.Find("td").Each(func(j int, selection *goquery.Selection) {
				if j == 4 {
					text := strings.ReplaceAll(selection.Text(), "\n", "")
					if text != "" {
						trust := "+" + text
						if i == 1 {
							operator.OP.HP += trust
						}
						if i == 2 {
							operator.OP.ATK += trust
						}
						if i == 3 {
							operator.OP.DEF += trust
						}
						if i == 4 {
							operator.OP.Res += trust
						}
					}
				}
			})
		})

		// 潜能
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "潜能提升" {
				selection.Parent().NextFilteredUntil(".nomobile", ".nodesktop").Each(func(j int, selection *goquery.Selection) {
					selection.Find("td").Each(func(k int, selection *goquery.Selection) {
						var potential Potential
						potential.Rank = k + 1
						potential.Desc = strings.ReplaceAll(selection.Text(), "\n", "")
						operator.Potentials = append(operator.Potentials, potential)
					})
				})
			}
		})
		// 天赋
		var talents []Talent
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "天赋" {
				selection.Parent().NextFilteredUntil(".wikitable", ".nodesktop").Each(func(j int, selection *goquery.Selection) {
					selection.Find("td").Each(func(k int, selection *goquery.Selection) {
						if k%3 == 0 {
							if selection.Nodes[0].FirstChild.Data == "ul" {
								return
							}
							var talent Talent
							talentName, _ := selection.Html()
							talent.Evolve = strings.ReplaceAll(selection.Next().Text(), "\n", "")
							if talent.Evolve == "" {
								return
							}
							desc, _ := selection.Next().Next().Html()
							talent.Name = template.HTML(talentName)
							talent.Desc = template.HTML(desc)
							if talent.Evolve == "" && len(talents) > 1 && talents[len(talents)-1].Name == talent.Name {
								return
							}
							talents = append(talents, talent)
						}
					})
				})
			}
		})
		operator.Talents = talents

		// 基建技能
		var buildingSkills []BuildingSkill
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "后勤技能" {
				selection.Parent().NextFilteredUntil(".wikitable", ".mw-heading").Each(func(j int, selection *goquery.Selection) {
					selection.Find("td").Each(func(k int, selection *goquery.Selection) {
						var buildingSkill BuildingSkill
						if k%5 == 0 {
							buildingSkill.Evolve = selection.Text()
							img, _ := selection.Next().Children().Attr("src")
							if img != "" {
								buildingSkill.Icon = "https:" + img
							}
							buildingSkill.Name = selection.Next().Next().Text()
							buildingSkill.Desc = strings.ReplaceAll(selection.Next().Next().Next().Next().Text(), "\n", "")
							buildingSkills = append(buildingSkills, buildingSkill)
						}
					})
				})
			}
		})
		operator.BuildingSkills = buildingSkills

		// 技能
		var skills []Skill
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "技能" {
				selection.Parent().NextFilteredUntil(".nomobile ", ".mw-heading").Each(func(j int, selection *goquery.Selection) {
					var skill Skill
					selection.Find("tr").Eq(0).Find("td").Each(func(k int, selection *goquery.Selection) {
						if k == 0 {
							icon, _ := selection.Children().Children().Children().Children().Attr("src")
							skill.Icon = icon
						}
						if k == 1 {
							skill.Name = strings.ReplaceAll(selection.Text(), "\n", "")
						}
						if k == 2 {
							selection.Children().Each(func(i int, selection *goquery.Selection) {
								spType, _ := selection.Html()
								skill.SpType = append(skill.SpType, template.HTML(spType))
							})
						}
						if k == 3 {
							skillRange, _ := selection.Children().Html()
							skill.SkillRange = template.HTML(skillRange)
						}
					})
					selection.Find("tr").Each(func(i int, selection *goquery.Selection) {
						tds := selection.Find("td")
						if len(tds.Nodes) == 5 {
							tds.Each(func(j int, selection *goquery.Selection) {
								text := strings.ReplaceAll(selection.Text(), "\n", "")
								if j == 1 {
									desc, _ := selection.Html()
									skill.Desc = template.HTML(desc)
								}
								if j == 2 {
									skill.SpInit = text
								}
								if j == 3 {
									skill.SpCost = text
								}
								if j == 4 {
									skill.Duration = formatDuration(text)
								}
							})
						}
					})
					skills = append(skills, skill)
				})
			}
		})
		operator.Skills = skills

		// 攻击范围
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "攻击范围" {
				selection.Parent().NextFilteredUntil(".nomobile ", ".nodesktop").Each(func(j int, selection *goquery.Selection) {
					if j == 0 {
						tds := selection.Find("td")
						td := tds.Eq(len(tds.Nodes) - 1)
						attackRange, _ := td.Children().Html()
						operator.AttackRange = buildRangeDoc(attackRange)
					}
				})
			}
		})
	}
	return operator
}

// formatDuration 持续时间统一处理：纯数字补上秒单位，非数字（如“无限”）原样展示
func formatDuration(text string) string {
	d := strings.TrimSpace(text)
	if d == "" {
		return ""
	}
	if _, err := strconv.ParseFloat(d, 64); err == nil {
		return d + "s"
	}
	return d
}

// rangeDocStyle 攻击范围片段在 iframe 内使用的样式
const rangeDocStyle = `html,body{margin:0;padding:0;background:transparent;overflow:hidden;}` +
	`table{border-collapse:separate;border-spacing:3px;margin:0;}` +
	`td{width:20px;height:20px;padding:0;background:rgba(255,255,255,.12);text-align:center;vertical-align:middle;line-height:0;}` +
	`td img{width:16px;height:16px;object-fit:contain;filter:brightness(0) invert(1);opacity:.9;}` +
	`img{max-width:24px;max-height:24px;object-fit:contain;}`

// buildRangeDoc wiki 的攻击范围片段结构不可控（可能带闭合标签、嵌套表格、外链图片），
// 放进 iframe 的 srcdoc 里渲染：既能统一外观，又不会破坏卡片本身的布局。
func buildRangeDoc(fragment string) string {
	return `<html><head><meta charset="utf-8" /><style>` + rangeDocStyle + `</style></head><body>` + fragment + `</body></html>`
}
