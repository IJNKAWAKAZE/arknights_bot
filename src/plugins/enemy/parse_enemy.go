package enemy

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"html/template"
	"net/http"
	"net/url"
	"strings"
)

type Enemy struct {
	Name       string        `json:"name"`       // 名字
	Pic        string        `json:"pic"`        // 图片
	Desc       string        `json:"desc"`       // 描述
	EnemyRace  string        `json:"enemyRace"`  // 种类
	EnemyLevel string        `json:"enemyLevel"` // 地位级别
	AttackType string        `json:"attackType"` // 攻击方式
	Motion     string        `json:"motion"`     // 行动方式
	Ability    template.HTML `json:"ability"`    // 能力
	Levels     []Level       `json:"level"`      // 级别信息
}

type Level struct {
	Desc       string        `json:"desc"`       // 描述
	AttackType string        `json:"attackType"` // 攻击方式
	Motion     string        `json:"motion"`     // 行动方式
	HpRecovery string        `json:"hpRecovery"` // 生命恢复速度
	HP         string        `json:"hp"`         // 生命值
	ATK        string        `json:"atk"`        // 攻击
	DEF        string        `json:"def"`        // 防御
	Res        string        `json:"res"`        // 法抗
	ATKRadius  string        `json:"ATKRadius"`  // 攻击半径
	Weight     string        `json:"weight"`     // 重量
	MoveSpeed  string        `json:"moveSpeed"`  // 移动速度
	Interval   string        `json:"interval"`   // 攻击间隔
	DamageRes  string        `json:"damageRes"`  // 损伤抵抗
	ElementRes string        `json:"elementRes"` // 元素抵抗
	Ridicule   string        `json:"ridicule"`   // 基础嘲讽等级
	Point      string        `json:"point"`      // 进点损失
	Abnormal   string        `json:"abnormal"`   // 异常抗性
	Skills     []Skill       `json:"skills"`     // 技能
	Talent     template.HTML `json:"talent"`     // 天赋&能力
}

type Skill struct {
	Name   string        `json:"name"`   // 名称
	SpInit template.HTML `json:"spInit"` // 初始sp
	SpCost template.HTML `json:"spCost"` // 所需sp
	Desc   template.HTML `json:"desc"`   // 描述
}

func ParseEnemy(name string) Enemy {
	var enemy Enemy
	api := viper.GetString("api.wiki")
	resp, _ := http.Get(api + url.PathEscape(name))
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	// 基本属性
	trs := doc.Find(".hlist").First().Find("tr")
	enemy.Name = strings.ReplaceAll(trs.Eq(0).Text(), "\n", "")
	pic, _ := trs.Eq(1).Find(".enemyicon a img").Attr("src")
	enemy.Pic = pic
	enemy.Desc = strings.ReplaceAll(trs.Eq(1).Find("td").Eq(1).Text(), "\n", "")
	td3 := trs.Eq(3).Find("td")
	enemy.EnemyRace = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")
	enemy.EnemyLevel = strings.ReplaceAll(td3.Eq(1).Text(), "\n", "")
	enemy.AttackType = strings.ReplaceAll(td3.Eq(2).Text(), "\n", "")
	enemy.Motion = strings.ReplaceAll(td3.Eq(3).Text(), "\n", "")
	ability, _ := trs.Eq(5).Html()
	enemy.Ability = template.HTML(ability)

	// 级别详情
	var levels []Level
	doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() == "级别0" || selection.Text() == "级别B" {
			var level Level
			selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
				if j == 0 {
					trs := selection.Find("tr")
					td3 := trs.Eq(3).Find("td")
					level.Desc = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")

					td5 := trs.Eq(6).Find("td")
					level.AttackType = strings.ReplaceAll(td5.Eq(1).Text(), "\n", "")
					level.Motion = strings.ReplaceAll(td5.Eq(3).Text(), "\n", "")
					level.HpRecovery = strings.ReplaceAll(td5.Eq(5).Text(), "\n", "")

					td6 := trs.Eq(6).Find("td")
					level.HP = strings.ReplaceAll(td6.Eq(0).Text(), "\n", "")
					level.ATK = strings.ReplaceAll(td6.Eq(1).Text(), "\n", "")
					level.DEF = strings.ReplaceAll(td6.Eq(2).Text(), "\n", "")
					level.Res = strings.ReplaceAll(td6.Eq(3).Text(), "\n", "")
					level.ATKRadius = strings.ReplaceAll(td6.Eq(4).Text(), "\n", "")
					level.Weight = strings.ReplaceAll(td6.Eq(5).Text(), "\n", "")

					td8 := trs.Eq(8).Find("td")
					level.MoveSpeed = strings.ReplaceAll(td8.Eq(0).Text(), "\n", "")
					level.Interval = strings.ReplaceAll(td8.Eq(1).Text(), "\n", "")
					level.DamageRes = strings.ReplaceAll(td8.Eq(2).Text(), "\n", "")
					level.ElementRes = strings.ReplaceAll(td8.Eq(3).Text(), "\n", "")
					level.Ridicule = strings.ReplaceAll(td8.Eq(4).Text(), "\n", "")
					level.Point = strings.ReplaceAll(td8.Eq(5).Text(), "\n", "")

					td9 := trs.Eq(9).Find("td")
					level.Abnormal = strings.ReplaceAll(td9.Eq(0).Text(), "\n", "")
				}
				if j == 1 {
					var skills []Skill
					selection.Find("tr").Each(func(k int, selection *goquery.Selection) {
						if !selection.Children().First().Is("th") {
							tds := selection.Find("td")
							// 技能
							if len(tds.Nodes) == 4 {
								var skill Skill
								skill.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
								spInit, _ := tds.Eq(1).Html()
								skill.SpInit = template.HTML(spInit)
								spCost, _ := tds.Eq(2).Html()
								skill.SpCost = template.HTML(spCost)
								desc, _ := tds.Eq(3).Html()
								skill.Desc = template.HTML(desc)
								skills = append(skills, skill)
							}
							// 天赋&能力
							if len(tds.Nodes) == 1 {
								talent, _ := tds.Html()
								level.Talent = template.HTML(talent)
							}
						}
					})
					level.Skills = skills
				}
			})
			levels = append(levels, level)
		}
		if selection.Text() == "级别1" || selection.Text() == "级别A" {
			var level Level
			selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
				if j == 0 {
					trs := selection.Find("tr")
					td3 := trs.Eq(3).Find("td")
					level.Desc = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")

					td5 := trs.Eq(6).Find("td")
					level.AttackType = strings.ReplaceAll(td5.Eq(1).Text(), "\n", "")
					level.Motion = strings.ReplaceAll(td5.Eq(3).Text(), "\n", "")
					level.HpRecovery = strings.ReplaceAll(td5.Eq(5).Text(), "\n", "")

					td6 := trs.Eq(6).Find("td")
					level.HP = strings.ReplaceAll(td6.Eq(0).Text(), "\n", "")
					level.ATK = strings.ReplaceAll(td6.Eq(1).Text(), "\n", "")
					level.DEF = strings.ReplaceAll(td6.Eq(2).Text(), "\n", "")
					level.Res = strings.ReplaceAll(td6.Eq(3).Text(), "\n", "")
					level.ATKRadius = strings.ReplaceAll(td6.Eq(4).Text(), "\n", "")
					level.Weight = strings.ReplaceAll(td6.Eq(5).Text(), "\n", "")

					td8 := trs.Eq(8).Find("td")
					level.MoveSpeed = strings.ReplaceAll(td8.Eq(0).Text(), "\n", "")
					level.Interval = strings.ReplaceAll(td8.Eq(1).Text(), "\n", "")
					level.DamageRes = strings.ReplaceAll(td8.Eq(2).Text(), "\n", "")
					level.ElementRes = strings.ReplaceAll(td8.Eq(3).Text(), "\n", "")
					level.Ridicule = strings.ReplaceAll(td8.Eq(4).Text(), "\n", "")
					level.Point = strings.ReplaceAll(td8.Eq(5).Text(), "\n", "")

					td9 := trs.Eq(9).Find("td")
					level.Abnormal = strings.ReplaceAll(td9.Eq(0).Text(), "\n", "")
				}
				if j == 1 {
					var skills []Skill
					selection.Find("tr").Each(func(k int, selection *goquery.Selection) {
						if !selection.Children().First().Is("th") {
							tds := selection.Find("td")
							// 技能
							if len(tds.Nodes) == 4 {
								var skill Skill
								skill.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
								spInit, _ := tds.Eq(1).Html()
								skill.SpInit = template.HTML(spInit)
								spCost, _ := tds.Eq(2).Html()
								skill.SpCost = template.HTML(spCost)
								desc, _ := tds.Eq(3).Html()
								skill.Desc = template.HTML(desc)
								skills = append(skills, skill)
							}
							// 天赋&能力
							if len(tds.Nodes) == 1 {
								talent, _ := tds.Html()
								level.Talent = template.HTML(talent)
							}
						}
					})
					level.Skills = skills
				}
			})
			levels = append(levels, level)
		}
		if selection.Text() == "级别2" || selection.Text() == "级别S" {
			var level Level
			selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
				if j == 0 {
					trs := selection.Find("tr")
					td3 := trs.Eq(3).Find("td")
					level.Desc = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")

					td5 := trs.Eq(6).Find("td")
					level.AttackType = strings.ReplaceAll(td5.Eq(1).Text(), "\n", "")
					level.Motion = strings.ReplaceAll(td5.Eq(3).Text(), "\n", "")
					level.HpRecovery = strings.ReplaceAll(td5.Eq(5).Text(), "\n", "")

					td6 := trs.Eq(6).Find("td")
					level.HP = strings.ReplaceAll(td6.Eq(0).Text(), "\n", "")
					level.ATK = strings.ReplaceAll(td6.Eq(1).Text(), "\n", "")
					level.DEF = strings.ReplaceAll(td6.Eq(2).Text(), "\n", "")
					level.Res = strings.ReplaceAll(td6.Eq(3).Text(), "\n", "")
					level.ATKRadius = strings.ReplaceAll(td6.Eq(4).Text(), "\n", "")
					level.Weight = strings.ReplaceAll(td6.Eq(5).Text(), "\n", "")

					td8 := trs.Eq(8).Find("td")
					level.MoveSpeed = strings.ReplaceAll(td8.Eq(0).Text(), "\n", "")
					level.Interval = strings.ReplaceAll(td8.Eq(1).Text(), "\n", "")
					level.DamageRes = strings.ReplaceAll(td8.Eq(2).Text(), "\n", "")
					level.ElementRes = strings.ReplaceAll(td8.Eq(3).Text(), "\n", "")
					level.Ridicule = strings.ReplaceAll(td8.Eq(4).Text(), "\n", "")
					level.Point = strings.ReplaceAll(td8.Eq(5).Text(), "\n", "")

					td9 := trs.Eq(9).Find("td")
					level.Abnormal = strings.ReplaceAll(td9.Eq(0).Text(), "\n", "")
				}
				if j == 1 {
					var skills []Skill
					selection.Find("tr").Each(func(k int, selection *goquery.Selection) {
						if !selection.Children().First().Is("th") {
							tds := selection.Find("td")
							// 技能
							if len(tds.Nodes) == 4 {
								var skill Skill
								skill.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
								spInit, _ := tds.Eq(1).Html()
								skill.SpInit = template.HTML(spInit)
								spCost, _ := tds.Eq(2).Html()
								skill.SpCost = template.HTML(spCost)
								desc, _ := tds.Eq(3).Html()
								skill.Desc = template.HTML(desc)
								skills = append(skills, skill)
							}
							// 天赋&能力
							if len(tds.Nodes) == 1 {
								talent, _ := tds.Html()
								level.Talent = template.HTML(talent)
							}
						}
					})
					level.Skills = skills
				}
			})
			levels = append(levels, level)
		}
	})
	enemy.Levels = levels
	return enemy
}
