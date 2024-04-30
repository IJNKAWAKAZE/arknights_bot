package web

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type CalendarInfo struct {
	Title string `json:"title"`
	Begin string `json:"begin"`
	Close string `json:"close"`
	End   string `json:"end"`
}

func Calendar(r *gin.Engine) {
	r.GET("/calendar", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Calendar.tmpl")
		var calendarMap = make(map[string]template.HTML)
		resp, err := http.Get(viper.GetString("api.calendar"))
		if err != nil {
			log.Println(err)
			return
		}
		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		text := doc.Text()
		begin := strings.Index(text, "{") + 1
		end := strings.Index(text, "}")
		reg := regexp.MustCompile("(\\[).*?(])")
		aaa := reg.FindAllStringSubmatch(text[begin:end], -1)
		var calendarInfo []CalendarInfo
		for _, a := range aaa {
			j := gjson.Parse(a[0])
			var c CalendarInfo
			c.Title = j.Get("0").String()
			c.Begin = fmt.Sprintf("%d-%02d-%02d", j.Get("1").Int(), j.Get("2").Int(), j.Get("3").Int())
			c.End = fmt.Sprintf("%d-%02d-%02d", j.Get("4").Int(), j.Get("5").Int(), j.Get("6").Int())
			if len(j.Array()) > 7 {
				c.Close = fmt.Sprintf("%d-%02d-%02d", j.Get("7").Int(), j.Get("8").Int(), j.Get("9").Int())
			}
			calendarInfo = append(calendarInfo, c)
		}
		defer resp.Body.Close()
		for _, c := range calendarInfo {
			beginTime, _ := time.ParseInLocation("2006-01-02", c.Begin, time.Local)
			//endTime, _ := time.ParseInLocation("2006-01-02", c.End, time.Local)
			//closeTime, _ := time.ParseInLocation("2006-01-02", c.Close, time.Local)
			if beginTime.Year() == time.Now().Year() {
				title := c.Title
				if _, bHas := calendarMap[c.Begin]; bHas {
					calendarMap[c.Begin] = template.HTML(fmt.Sprintf("%s<br>开始 %s", calendarMap[c.Begin], title))
				} else {
					calendarMap[c.Begin] = template.HTML("开始 " + title)
				}
				if _, eHas := calendarMap[c.End]; eHas {
					calendarMap[c.End] = template.HTML(fmt.Sprintf("%s<br>结束 %s", calendarMap[c.End], title))
				} else {
					calendarMap[c.End] = template.HTML("结束 " + c.Title)
				}
				if c.Close != "" {
					if _, cHas := calendarMap[c.Close]; cHas {
						calendarMap[c.Close] = template.HTML(fmt.Sprintf("%s<br>关闭关卡 %s", calendarMap[c.Close], title))
					} else {
						calendarMap[c.Close] = template.HTML("关闭关卡" + c.Title)
					}
				}
			}
		}
		c.HTML(http.StatusOK, "Calendar.tmpl", calendarMap)
	})
}
