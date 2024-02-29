package operator

import (
	"arknights_bot/utils"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"html/template"
	"net/http"
	"regexp"
	"strings"
)

type Operator struct {
	OP               utils.Operator   `json:"op"`
	Painting         string           `json:"painting"`
	AttackRange      template.HTML    `json:"attackRange"`
	ProfessionBranch ProfessionBranch `json:"professionBranch"`
	Potentials       []Potential      `json:"potentials"`
	Talents          []Talent         `json:"talents"`
	BuildingSkills   []BuildingSkill  `json:"buildingSkills"`
	Skills           []Skill          `json:"skills"`
}

type ProfessionBranch struct {
	Name string `json:"name"`
	Pic  string `json:"pic"`
	Desc string `json:"desc"`
}

type Potential struct {
	Rank int    `json:"rank"`
	Desc string `json:"desc"`
}

type Talent struct {
	Evolve string `json:"evolve"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
}

type BuildingSkill struct {
	Evolve string `json:"evolve"`
	Icon   string `json:"icon"`
	Name   string `json:"name"`
	Desc   string `json:"desc"`
}

type Skill struct {
	Icon       string          `json:"icon"`
	Name       string          `json:"name"`
	Desc       string          `json:"desc"`
	SkillRange template.HTML   `json:"skillRange"`
	SpType     []template.HTML `json:"spType"`
	SpInit     string          `json:"spInit"`
	SpCost     string          `json:"spCost"`
	Duration   string          `json:"duration"`
}

// ParseOperator 解析干员数据
func ParseOperator(name string) Operator {
	var operator Operator
	api := viper.GetString("api.wiki")
	response, _ := http.Get(api + name)
	op := utils.GetOperatorByName(name)
	if op.Name != "" {
		operator.OP = op
		operator.Painting = op.Skins[0]
		doc, _ := goquery.NewDocumentFromReader(response.Body)

		// 职业分支
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "特性" {
				selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
					tds := selection.Find("td")
					operator.ProfessionBranch.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
					paintingName := fmt.Sprintf("职业分支图标_%s.png", operator.ProfessionBranch.Name)
					m := utils.Md5(paintingName)
					path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
					operator.ProfessionBranch.Pic = path + paintingName
					tds.Each(func(j int, selection *goquery.Selection) {
						if _, b := selection.Attr("style"); !b {
							operator.ProfessionBranch.Desc = strings.ReplaceAll(selection.Text(), "\n", "")
							fmt.Println(j, "-", selection.Text())
						}
					})
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
				selection.NextFilteredUntil(".nomobile", "h2").Each(func(j int, selection *goquery.Selection) {
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
				selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
					selection.Find("td").Each(func(k int, selection *goquery.Selection) {
						if k%3 == 0 {
							var talent Talent
							talentName := strings.ReplaceAll(selection.Text(), "\n", "")
							talent.Evolve = strings.ReplaceAll(selection.Next().Text(), "\n", "")
							desc := strings.ReplaceAll(selection.Next().Next().Text(), "\n", "")
							reg := regexp.MustCompile("(\\[).*?(])")
							talent.Name = talentName
							talent.Desc = reg.ReplaceAllString(desc, "")
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
				selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
					var buildingSkill BuildingSkill
					buildingSkill.Evolve = selection.Find("td").Eq(0).Text()
					img, _ := selection.Find("td").Eq(1).Children().Attr("data-src")
					buildingSkill.Icon = "https:" + img
					buildingSkill.Name = selection.Find("td").Eq(2).Text()
					buildingSkill.Desc = strings.ReplaceAll(selection.Find("td").Eq(4).Text(), "\n", "")
					buildingSkills = append(buildingSkills, buildingSkill)
				})
			}
		})
		operator.BuildingSkills = buildingSkills

		// 技能
		var skills []Skill
		doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
			if selection.Text() == "技能" {
				selection.NextFilteredUntil(".nomobile ", "h2").Each(func(j int, selection *goquery.Selection) {
					var skill Skill
					selection.Find("tr").Eq(0).Find("td").Each(func(k int, selection *goquery.Selection) {
						if k == 0 {
							icon, _ := selection.Children().Children().Children().Attr("data-src")
							skill.Icon = "https://prts.wiki" + icon
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
									reg := regexp.MustCompile("(\\[).*?(])")
									skill.Desc = reg.ReplaceAllString(text, "")
								}
								if j == 2 {
									skill.SpInit = text
								}
								if j == 3 {
									skill.SpCost = text
								}
								if j == 4 {
									skill.Duration = text
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
				selection.NextFilteredUntil(".nomobile ", "h2").Each(func(j int, selection *goquery.Selection) {
					tds := selection.Find("td")
					td := tds.Eq(len(tds.Nodes) - 1)
					attackRange, _ := td.Children().Html()
					operator.AttackRange = template.HTML(attackRange)
					return
				})
			}
		})
	}
	return operator
}
