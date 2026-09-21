package operator

import (
	"arknights_bot/utils/hashutil"
	"arknights_bot/utils/httpx"
	"arknights_bot/utils/localassets"
	"arknights_bot/utils/model"
	"arknights_bot/utils/search"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"html/template"
	"math"
	"regexp"
	"strconv"
	"strings"
)

type Operator struct {
	OP               model.Operator   `json:"op"`               // 基本信息
	Painting         string           `json:"painting"`         // 立绘
	AttackRanges     []RangeGroup     `json:"attackRanges"`     // 攻击范围（按精英化阶段）
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
		doc, _ := goquery.NewDocumentFromReader(httpx.Body(api + name))

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
					operator.ProfessionBranch.Pic = path + localassets.EscapePath(paintingName)
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
						operator.AttackRanges = parseAttackRanges(selection)
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

// 攻击范围 SVG 的坐标步长：22px 的格子在 viewBox 里占 26 个单位（含 4 单位间隔）。
const (
	rangeCellUnit = 26.0  // viewBox 中一个格子的步长
	rangeCellSize = 18.0  // 卡片上格子的显示边长
	rangeMaxSide  = 108.0 // 单张范围图的最大边长，兜底防止异常大的范围撑破卡片
)

var (
	svgTagPattern    = regexp.MustCompile(`(?is)<svg\b[^>]*>`)
	viewBoxPattern   = regexp.MustCompile(`(?i)viewbox\s*=\s*"([^"]*)"`)
	styleAttrPattern = regexp.MustCompile(`(?is)\sstyle\s*=\s*"([^"]*)"`)
)

// RangeGroup 攻击范围按精英化阶段分组，精英0/精英1/精英2 各一张范围图
type RangeGroup struct {
	Label string        `json:"label"` // 精英化阶段
	Grid  template.HTML `json:"grid"`  // 范围图（已按卡片尺寸缩放）
}

// parseAttackRanges wiki 的范围表格里标签与范围图是分开排布的
// （桌面版一行表头加一行范围图，移动版每行一组），按出现顺序配对即可兼容两种版式。
func parseAttackRanges(table *goquery.Selection) []RangeGroup {
	if table == nil {
		return nil
	}
	var labels []string
	table.Find("th").Each(func(_ int, th *goquery.Selection) {
		label := strings.Join(strings.Fields(th.Text()), "")
		span, err := strconv.Atoi(th.AttrOr("colspan", "1"))
		if err != nil || span < 1 {
			span = 1
		}
		for i := 0; i < span; i++ {
			labels = append(labels, label)
		}
	})

	var groups []RangeGroup
	table.Find("svg,img").Each(func(i int, grid *goquery.Selection) {
		html, err := goquery.OuterHtml(grid)
		if err != nil {
			return
		}
		group := RangeGroup{Grid: template.HTML(scaleRangeGrid(html))}
		if i < len(labels) {
			group.Label = labels[i]
		}
		groups = append(groups, group)
	})
	return groups
}

// scaleRangeGrid wiki 给的范围图带 130px 的固定宽高，精英2 这类大范围会超出卡片被裁掉；
// 这里统一缩放到固定格子大小，不同大小的范围格子一致，也不会再出现截断。
func scaleRangeGrid(html string) string {
	return svgTagPattern.ReplaceAllStringFunc(html, func(tag string) string {
		width, height := svgViewBoxSize(tag)
		if width <= 0 || height <= 0 {
			return tag
		}
		maxSide := math.Max(width, height)
		scale := rangeCellSize / rangeCellUnit
		if side := maxSide * scale; side > rangeMaxSide {
			scale = rangeMaxSide / maxSide
		}
		return setInlineSize(tag, width*scale, height*scale)
	})
}

// svgViewBoxSize 读取 svg 的 viewBox 宽高，用于推算格子数量与缩放比例
func svgViewBoxSize(tag string) (float64, float64) {
	match := viewBoxPattern.FindStringSubmatch(tag)
	if match == nil {
		return 0, 0
	}
	fields := strings.FieldsFunc(match[1], func(r rune) bool {
		return r == ' ' || r == ',' || r == '\t' || r == '\n' || r == '\r'
	})
	if len(fields) != 4 {
		return 0, 0
	}
	width, errW := strconv.ParseFloat(fields[2], 64)
	height, errH := strconv.ParseFloat(fields[3], 64)
	if errW != nil || errH != nil || width <= 0 || height <= 0 {
		return 0, 0
	}
	return width, height
}

// setInlineSize 把尺寸追加进 style（写在原有声明之后，覆盖 wiki 自带的固定宽高）
func setInlineSize(tag string, width, height float64) string {
	size := fmt.Sprintf("width:%.1fpx;height:%.1fpx;", width, height)
	if loc := styleAttrPattern.FindStringSubmatchIndex(tag); loc != nil {
		// wiki 的样式以 width:!important 这种半截声明结尾，直接拼接会和前面的声明粘在一起被整条丢弃
		sep := ";"
		if existing := strings.TrimSpace(tag[loc[2]:loc[3]]); existing == "" || strings.HasSuffix(existing, ";") {
			sep = ""
		}
		return tag[:loc[3]] + sep + size + tag[loc[3]:]
	}
	return strings.TrimSuffix(tag, ">") + ` style="` + size + `">`
}
