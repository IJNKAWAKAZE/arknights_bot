package arknightsnews

import (
	"arknights_bot/config"
	"arknights_bot/utils"
	"bytes"
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"io"
	"net/http"
	"os"
	"strings"
)

type Payload struct {
	Payload string `json:"payload"`
}

type Pic struct {
	Url    string `json:"url"`
	Height int64  `json:"height"`
}

func BilibiliNews() {
	text, pics := ParseBilibiliDynamic()
	if len(text) == 0 {
		return
	}
	groups := utils.GetJoinedGroups()
	if pics == nil {
		for _, group := range groups {
			sendMessage := tgbotapi.NewMessage(group, text)
			config.Arknights.Send(sendMessage)
		}
		return
	}

	if len(pics) == 1 {
		for _, group := range groups {
			if pics[0].Height > 3000 {
				sendDocument := tgbotapi.NewDocument(group, tgbotapi.FileURL(pics[0].Url))
				sendDocument.Caption = text
				config.Arknights.Send(sendDocument)
			} else {
				sendPhoto := tgbotapi.NewPhoto(group, tgbotapi.FileURL(pics[0].Url))
				sendPhoto.Caption = text
				config.Arknights.Send(sendPhoto)
			}
		}
		return
	}

	for _, group := range groups {
		var mediaGroup tgbotapi.MediaGroupConfig
		var media []interface{}
		mediaGroup.ChatID = group

		d := false
		for _, p := range pics {
			if p.Height > 3000 {
				d = true
			}
		}

		for i, pic := range pics {
			if d {
				var inputDocument tgbotapi.InputMediaDocument
				inputDocument.Media = tgbotapi.FileURL(pic.Url)
				inputDocument.Type = "document"
				if i == len(pics)-1 {
					inputDocument.Caption = text
				}
				media = append(media, inputDocument)
			} else {
				if strings.HasSuffix(pic.Url, ".gif") {
					var inputVideo tgbotapi.InputMediaVideo
					inputVideo.Media = tgbotapi.FileBytes{Bytes: convert2Video(pic.Url, i)}
					inputVideo.Type = "video"
					media = append(media, inputVideo)
					continue
				}
				var inputPhoto tgbotapi.InputMediaPhoto
				inputPhoto.Media = tgbotapi.FileURL(pic.Url)
				inputPhoto.Type = "photo"
				if i == 0 {
					inputPhoto.Caption = text
				}
				media = append(media, inputPhoto)
			}
		}
		mediaGroup.Media = media
		config.Arknights.SendMediaGroup(mediaGroup)
	}
}

func convert2Video(url string, i int) []byte {
	outPut := fmt.Sprintf("./temp%d.mp4", i)
	res, _ := http.Get(url)
	tempFile, _ := os.CreateTemp("./", "temp-*.gif")
	io.Copy(tempFile, res.Body)
	tempFile.Close()
	defer res.Body.Close()
	ffmpeg.Input(tempFile.Name()).
		Output(outPut).
		OverWriteOutput().Run()
	os.Remove(tempFile.Name())
	f, _ := os.Open(outPut)
	b, _ := io.ReadAll(f)
	f.Close()
	os.Remove(f.Name())
	return b
}

func ParseBilibiliDynamic() (string, []Pic) {
	var text string
	var pics []Pic
	b3, b4, err := generateBuvid()
	if err != nil {
		return text, pics
	}
	err = registerBuvid(b3, b4)
	if err != nil {
		return text, pics
	}
	url := viper.GetString("api.bilibili_dynamic")
	resBody, err := requestBili("GET", fmt.Sprintf("buvid3=%s; buvid4=%s", b3, b4), url, nil)
	if err != nil {
		return text, pics
	}
	result := gjson.ParseBytes(resBody)
	items := result.Get("data.items").Array()
	for _, item := range items {
		top := item.Get("modules.module_tag.text").String()
		if top != "置顶" {
			dynamicType := item.Get("type").String()
			id := item.Get("id_str").String()
			link := "https://t.bilibili.com/" + id
			//publishTime := time.Unix(item.Get("modules.module_author.pub_ts").Int(), 0).Format("2006-01-02 15:04:05")
			if dynamicType == "DYNAMIC_TYPE_DRAW" {
				for _, pic := range item.Get("modules.module_dynamic.major.opus.pics").Array() {
					var p Pic
					p.Url = pic.Get("url").String()
					p.Height = pic.Get("height").Int()
					pics = append(pics, p)
				}
				text = item.Get("modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_WORD" {
				text = item.Get("modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_AV" {
				title := item.Get("modules.module_dynamic.major.archive.title").String() + "\n\n"
				desc := item.Get("modules.module_dynamic.major.archive.desc").String()
				cover := item.Get("modules.module_dynamic.major.archive.cover").String()
				vUrl := "https:" + item.Get("modules.module_dynamic.major.archive.jump_url").String()
				text = title + desc + "\n视频链接：" + vUrl
				var p Pic
				p.Url = cover
				pics = append(pics, p)
			}
			if dynamicType == "DYNAMIC_TYPE_FORWARD" {
				desc := item.Get("modules.module_dynamic.desc.text").String()
				for _, pic := range item.Get("orig.modules.module_dynamic.major.opus.pics").Array() {
					var p Pic
					p.Url = pic.Get("url").String()
					p.Height = pic.Get("height").Int()
					pics = append(pics, p)
				}
				text = desc + "\n\n" + item.Get("orig.modules.module_dynamic.major.opus.summary.text").String()
			}
			if dynamicType == "DYNAMIC_TYPE_ARTICLE" {
				summary := item.Get("modules.module_dynamic.major.opus.summary.text").String()
				for _, pic := range item.Get("modules.module_dynamic.major.opus.pics").Array() {
					var p Pic
					p.Url = pic.Get("url").String()
					p.Height = pic.Get("height").Int()
					pics = append(pics, p)
				}
				text = strings.ReplaceAll(summary, "[图片]", "") + "\n\n专栏地址：https:" + item.Get("modules.module_dynamic.major.opus.jump_url").String()
			}
			if utils.RedisSetIsExists("tg_arknights", link) {
				return "", nil
			}
			utils.RedisAddSet("tg_arknights", link)
			break
		}
	}
	return text, pics
}

func generateBuvid() (string, string, error) {
	url := viper.GetString("api.bilibili_buvid")
	resBody, err := requestBili("GET", "", url, nil)
	if err != nil {
		return "", "", err
	}
	jsonData := gjson.ParseBytes(resBody)
	b3 := jsonData.Get("data.b_3").String()
	b4 := jsonData.Get("data.b_4").String()
	return b3, b4, nil
}

func registerBuvid(b3, b4 string) error {
	url := viper.GetString("api.bilibili_register_buvid")
	jsonData := `{"3064":2,"5062":"1704899411253","03bf":"","39c8":"333.937.fp.risk","34f1":"","d402":"","654a":"","6e7c":"360x668","3c43":{"2673":0,"5766":24,"6527":0,"7003":1,"807e":1,"b8ce":"Mozilla/5.0 (Linux; Android 10; K) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Mobile Safari/537.36 EdgA/118.0.2088.66","641c":0,"07a4":"zh-CN","1c57":4,"0bd0":8,"fc9d":-480,"6aa9":"Asia/Shanghai","75b8":1,"3b21":1,"8a1c":1,"d52f":"not available","adca":"Linux armv81","80c9":[],"13ab":"zMgAAAAASUVORK5CYII=","bfe9":"mgQDEKAKxirCZRFLCvwP8Bjez5pveZop4AAAAASUVORK5CYII=","6bc5":"Google Inc. (ARM)~ANGLE (ARM, Mali-G57 MC2, OpenGL ES 3.2)","ed31":0,"72bd":0,"097b":0,"d02f":"124.08072766105033"},"54ef":"{}","8b94":"","df35":"A95D3545-DEC10-D817-35410-531784C2281905903infoc","07a4":"zh-CN","5f45":null,"db46":0}`
	payload := Payload{
		Payload: jsonData,
	}
	payloadb, _ := json.Marshal(payload)
	_, err := requestBili("POST", fmt.Sprintf("buvid3=%s; buvid4=%s", b3, b4), url, bytes.NewReader(payloadb))
	if err != nil {
		return err
	}
	return nil
}

func requestBili(method, cookie, url string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("User-Agent", viper.GetString("api.user_agent"))
	req.Header.Add("referer", "https://m.bilibili.com/")
	req.Header.Add("Content-Type", "application/json")
	if cookie != "" {
		req.Header.Add("Cookie", cookie)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	resBody, _ := io.ReadAll(res.Body)
	defer res.Body.Close()
	return resBody, nil
}
