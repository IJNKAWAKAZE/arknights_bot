package enemy

import (
	"arknights_bot/utils"
	"fmt"
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
	Endure     string        `json:"endure"`     // 耐久
	Attack     string        `json:"attack"`     // 攻击力
	Defence    string        `json:"defence"`    // 防御力
	Resistance string        `json:"resistance"` // 法抗
	Ability    template.HTML `json:"ability"`    // 能力
	Levels     []Level       `json:"level"`      // 级别信息
}

type Level struct {
	HP         string        `json:"hp"`         // 生命值
	ATK        string        `json:"atk"`        // 攻击
	DEF        string        `json:"def"`        // 防御
	Res        string        `json:"res"`        // 法抗
	Interval   string        `json:"interval"`   // 攻击间隔
	ATKRadius  string        `json:"ATKRadius"`  // 攻击半径
	Point      string        `json:"point"`      // 进点损失
	MoveSpeed  string        `json:"moveSpeed"`  // 移动速度
	Weight     string        `json:"weight"`     // 重量
	ElementRes string        `json:"elementRes"` // 元素抵抗
	DamageRes  string        `json:"damageRes"`  // 损伤抵抗
	HpRecovery string        `json:"hpRecovery"` // 生命恢复速度
	Silence    string        `json:"silence"`    // 沉默免疫
	Dizziness  string        `json:"dizziness"`  // 眩晕免疫
	Sleep      string        `json:"sleep"`      // 睡眠免疫
	Freeze     string        `json:"freeze"`     // 冻结免疫
	Float      string        `json:"float"`      // 浮空免疫
	Tremble    string        `json:"tremble"`    // 战栗免疫
	Skills     []Skill       `json:"skills"`     // 技能
	Talent     template.HTML `json:"talent"`     // 天赋&能力
}

type Skill struct {
	Name      string        `json:"name"`      // 名称
	SpInit    string        `json:"spInit"`    // 初始sp
	SpCost    string        `json:"spCost"`    // 所需sp
	SkillType template.HTML `json:"skillType"` // 技能类型
	Desc      template.HTML `json:"desc"`      // 描述
}

func ParseEnemy(name string) Enemy {
	var enemy Enemy
	api := viper.GetString("api.wiki")
	resp, _ := http.Get(api + url.PathEscape(name))
	doc, _ := goquery.NewDocumentFromReader(resp.Body)
	// 基本属性
	trs := doc.Find(".hlist").First().Find("tr")
	enemy.Name = strings.ReplaceAll(trs.Eq(0).Text(), "\n", "")
	paintingName := fmt.Sprintf("头像_敌人_%s.png", name)
	m := utils.Md5(paintingName)
	path := "https://prts.wiki" + fmt.Sprintf("/images/%s/%s/", m[:1], m[:2])
	enemy.Pic = path + paintingName
	enemy.Desc = strings.ReplaceAll(trs.Eq(1).Find("td").Eq(1).Text(), "\n", "")
	td3 := trs.Eq(3).Find("td")
	enemy.EnemyRace = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")
	enemy.EnemyLevel = strings.ReplaceAll(td3.Eq(1).Text(), "\n", "")
	enemy.AttackType = strings.ReplaceAll(td3.Eq(2).Text(), "\n", "")
	enemy.Motion = strings.ReplaceAll(td3.Eq(3).Text(), "\n", "")
	td5 := trs.Eq(5).Find("td")
	enemy.Endure = strings.ReplaceAll(td5.Eq(0).Text(), "\n", "")
	enemy.Attack = strings.ReplaceAll(td5.Eq(1).Text(), "\n", "")
	enemy.Defence = strings.ReplaceAll(td5.Eq(2).Text(), "\n", "")
	enemy.Resistance = strings.ReplaceAll(td5.Eq(3).Text(), "\n", "")
	ability, _ := trs.Eq(7).Html()
	enemy.Ability = template.HTML(ability)

	// 级别详情
	var levels []Level
	doc.Find("h2").Each(func(i int, selection *goquery.Selection) {
		if selection.Text() == "级别0" {
			var level Level
			selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
				if j == 0 {
					trs := selection.Find("tr")
					td3 := trs.Eq(3).Find("td")
					level.HP = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")
					level.ATK = strings.ReplaceAll(td3.Eq(1).Text(), "\n", "")
					level.DEF = strings.ReplaceAll(td3.Eq(2).Text(), "\n", "")
					level.Res = strings.ReplaceAll(td3.Eq(3).Text(), "\n", "")
					level.Interval = strings.ReplaceAll(td3.Eq(4).Text(), "\n", "")
					level.ATKRadius = strings.ReplaceAll(td3.Eq(5).Text(), "\n", "")

					td5 := trs.Eq(5).Find("td")
					level.Point = strings.ReplaceAll(td5.Eq(0).Text(), "\n", "")
					level.MoveSpeed = strings.ReplaceAll(td5.Eq(1).Text(), "\n", "")
					level.Weight = strings.ReplaceAll(td5.Eq(2).Text(), "\n", "")
					level.ElementRes = strings.ReplaceAll(td5.Eq(3).Text(), "\n", "")
					level.DamageRes = strings.ReplaceAll(td5.Eq(4).Text(), "\n", "")
					level.HpRecovery = strings.ReplaceAll(td5.Eq(5).Text(), "\n", "")

					td8 := trs.Eq(8).Find("td")
					level.Silence = strings.ReplaceAll(td8.Eq(0).Text(), "\n", "")
					level.Dizziness = strings.ReplaceAll(td8.Eq(1).Text(), "\n", "")
					level.Sleep = strings.ReplaceAll(td8.Eq(2).Text(), "\n", "")
					level.Freeze = strings.ReplaceAll(td8.Eq(3).Text(), "\n", "")
					level.Float = strings.ReplaceAll(td8.Eq(4).Text(), "\n", "")
					level.Tremble = strings.ReplaceAll(td8.Eq(5).Text(), "\n", "")
				}
				if j == 1 {
					var skills []Skill
					selection.Find("tr").Each(func(k int, selection *goquery.Selection) {
						if !selection.Children().First().Is("th") {
							tds := selection.Find("td")
							// 技能
							if len(tds.Nodes) == 5 {
								var skill Skill
								skill.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
								skill.SpInit = strings.ReplaceAll(tds.Eq(1).Text(), "\n", "")
								skill.SpCost = strings.ReplaceAll(tds.Eq(2).Text(), "\n", "")
								skillType, _ := tds.Eq(3).Html()
								skill.SkillType = template.HTML(skillType)
								desc, _ := tds.Eq(4).Html()
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
		if selection.Text() == "级别1" {
			var level Level
			selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
				if j == 0 {
					trs := selection.Find("tr")
					td3 := trs.Eq(3).Find("td")
					level.HP = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")
					level.ATK = strings.ReplaceAll(td3.Eq(1).Text(), "\n", "")
					level.DEF = strings.ReplaceAll(td3.Eq(2).Text(), "\n", "")
					level.Res = strings.ReplaceAll(td3.Eq(3).Text(), "\n", "")
					level.Interval = strings.ReplaceAll(td3.Eq(4).Text(), "\n", "")
					level.ATKRadius = strings.ReplaceAll(td3.Eq(5).Text(), "\n", "")

					td5 := trs.Eq(5).Find("td")
					level.Point = strings.ReplaceAll(td5.Eq(0).Text(), "\n", "")
					level.MoveSpeed = strings.ReplaceAll(td5.Eq(1).Text(), "\n", "")
					level.Weight = strings.ReplaceAll(td5.Eq(2).Text(), "\n", "")
					level.ElementRes = strings.ReplaceAll(td5.Eq(3).Text(), "\n", "")
					level.DamageRes = strings.ReplaceAll(td5.Eq(4).Text(), "\n", "")
					level.HpRecovery = strings.ReplaceAll(td5.Eq(5).Text(), "\n", "")

					td8 := trs.Eq(8).Find("td")
					level.Silence = strings.ReplaceAll(td8.Eq(0).Text(), "\n", "")
					level.Dizziness = strings.ReplaceAll(td8.Eq(1).Text(), "\n", "")
					level.Sleep = strings.ReplaceAll(td8.Eq(2).Text(), "\n", "")
					level.Freeze = strings.ReplaceAll(td8.Eq(3).Text(), "\n", "")
					level.Float = strings.ReplaceAll(td8.Eq(4).Text(), "\n", "")
					level.Tremble = strings.ReplaceAll(td8.Eq(5).Text(), "\n", "")
				}
				if j == 1 {
					var skills []Skill
					selection.Find("tr").Each(func(k int, selection *goquery.Selection) {
						if !selection.Children().First().Is("th") {
							tds := selection.Find("td")
							// 技能
							if len(tds.Nodes) == 5 {
								var skill Skill
								skill.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
								skill.SpInit = strings.ReplaceAll(tds.Eq(1).Text(), "\n", "")
								skill.SpCost = strings.ReplaceAll(tds.Eq(2).Text(), "\n", "")
								skillType, _ := tds.Eq(3).Html()
								skill.SkillType = template.HTML(skillType)
								desc, _ := tds.Eq(4).Html()
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
		if selection.Text() == "级别2" {
			var level Level
			selection.NextFilteredUntil(".wikitable", "h2").Each(func(j int, selection *goquery.Selection) {
				if j == 0 {
					trs := selection.Find("tr")
					td3 := trs.Eq(3).Find("td")
					level.HP = strings.ReplaceAll(td3.Eq(0).Text(), "\n", "")
					level.ATK = strings.ReplaceAll(td3.Eq(1).Text(), "\n", "")
					level.DEF = strings.ReplaceAll(td3.Eq(2).Text(), "\n", "")
					level.Res = strings.ReplaceAll(td3.Eq(3).Text(), "\n", "")
					level.Interval = strings.ReplaceAll(td3.Eq(4).Text(), "\n", "")
					level.ATKRadius = strings.ReplaceAll(td3.Eq(5).Text(), "\n", "")

					td5 := trs.Eq(5).Find("td")
					level.Point = strings.ReplaceAll(td5.Eq(0).Text(), "\n", "")
					level.MoveSpeed = strings.ReplaceAll(td5.Eq(1).Text(), "\n", "")
					level.Weight = strings.ReplaceAll(td5.Eq(2).Text(), "\n", "")
					level.ElementRes = strings.ReplaceAll(td5.Eq(3).Text(), "\n", "")
					level.DamageRes = strings.ReplaceAll(td5.Eq(4).Text(), "\n", "")
					level.HpRecovery = strings.ReplaceAll(td5.Eq(5).Text(), "\n", "")

					td8 := trs.Eq(8).Find("td")
					level.Silence = strings.ReplaceAll(td8.Eq(0).Text(), "\n", "")
					level.Dizziness = strings.ReplaceAll(td8.Eq(1).Text(), "\n", "")
					level.Sleep = strings.ReplaceAll(td8.Eq(2).Text(), "\n", "")
					level.Freeze = strings.ReplaceAll(td8.Eq(3).Text(), "\n", "")
					level.Float = strings.ReplaceAll(td8.Eq(4).Text(), "\n", "")
					level.Tremble = strings.ReplaceAll(td8.Eq(5).Text(), "\n", "")
				}
				if j == 1 {
					var skills []Skill
					selection.Find("tr").Each(func(k int, selection *goquery.Selection) {
						if !selection.Children().First().Is("th") {
							tds := selection.Find("td")
							// 技能
							if len(tds.Nodes) == 5 {
								var skill Skill
								skill.Name = strings.ReplaceAll(tds.Eq(0).Text(), "\n", "")
								skill.SpInit = strings.ReplaceAll(tds.Eq(1).Text(), "\n", "")
								skill.SpCost = strings.ReplaceAll(tds.Eq(2).Text(), "\n", "")
								skillType, _ := tds.Eq(3).Html()
								skill.SkillType = template.HTML(skillType)
								desc, _ := tds.Eq(4).Html()
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
