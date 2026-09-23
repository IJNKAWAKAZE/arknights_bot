package datasource

import (
	"arknights_bot/config"
	"arknights_bot/utils/cache"
	"arknights_bot/utils/hashutil"
	"arknights_bot/utils/httpx"
	"arknights_bot/utils/localassets"
	"arknights_bot/utils/model"
	"arknights_bot/utils/search"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/spf13/viper"
	"github.com/starudream/go-lib/core/v2/codec/json"
	"github.com/tidwall/gjson"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

var Profession = make(map[string]string)

func init() {
	profession := make(map[string]string)
	profession["术师"] = "CASTER"
	profession["医疗"] = "MEDIC"
	profession["先锋"] = "PIONEER"
	profession["狙击"] = "SNIPER"
	profession["特种"] = "SPECIAL"
	profession["辅助"] = "SUPPORT"
	profession["重装"] = "TANK"
	profession["近卫"] = "WARRIOR"
	Profession = profession
}

// UpdateDataSource 更新数据源
func UpdateDataSource() {
	if err := UpdateDataSourceRunner(); err != nil {
		log.Println("数据源更新失败:", err)
	}
}

// UpdateDataSourceRunner 更新数据源。
// 抓取失败或解析结果明显异常时返回错误并中止，绝不写入不完整的数据。
func UpdateDataSourceRunner() error {
	log.Println("开始更新数据源...")
	var operators []model.Operator
	api := viper.GetString("api.wiki")
	doc, ok := fetchDocument(api + "干员一览")
	if !ok {
		return fmt.Errorf("获取干员一览失败，本次更新已中止")
	}
	doc.Find("#filter-data div").Each(func(i int, selection *goquery.Selection) {
		var operator model.Operator
		attrs := selection.Nodes[0].Attr
		operator.Name = attrs[0].Val
		operator.Profession = Profession[attrs[1].Val]
		operator.ProfessionZH = attrs[1].Val + "干员"
		operator.Rarity, _ = strconv.Atoi(attrs[2].Val)
		operator.Logo = attrs[3].Val
		operator.Race = attrs[6].Val
		operator.NameEn = attrs[7].Val
		operator.NameJa = attrs[8].Val
		operator.Code = attrs[9].Val
		operator.HP = attrs[10].Val
		operator.ATK = attrs[11].Val
		operator.DEF = attrs[12].Val
		operator.Res = attrs[13].Val
		operator.ReDeploy = attrs[14].Val
		c := strings.Split(attrs[15].Val, "→")
		operator.Cost = c[len(c)-1]
		b := strings.Split(attrs[16].Val, "→")
		operator.Block = b[len(b)-1]
		operator.Interval = attrs[17].Val
		operator.Sex = attrs[18].Val
		operator.Position = attrs[19].Val
		operator.Tags = strings.ReplaceAll(attrs[20].Val, "支援机械", "机械")
		operator.ObtainMethod = attrs[21].Val
		// 头像
		paintingName := fmt.Sprintf("头像_%s.png", operator.Name)
		m := hashutil.Md5(paintingName)
		path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
		operator.Avatar = path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
		// 半身像
		paintingName = fmt.Sprintf("半身像_%s_1.png", operator.Name)
		m = hashutil.Md5(paintingName)
		path = "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
		operator.ThumbURL = path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
		operators = append(operators, operator)
	})

	// 页面结构变化或拿到错误页时解析结果会异常偏少，这时直接中止，避免用空数据覆盖已有缓存
	if len(operators) < minOperators {
		return fmt.Errorf("干员一览解析结果异常，仅解析到 %d 个干员，本次更新已中止", len(operators))
	}

	// 已有的生日缓存：既用来判断哪些干员还没记录到生日，也用来把这次重新查到的结果并回去
	birthdays, birthdayBuckets := loadBirthdayCache()

	skinCount := make(map[string][]string)
	skinDoc, ok := fetchDocument(api + "时装回廊")
	if !ok {
		return fmt.Errorf("获取时装回廊失败，本次更新已中止")
	}
	skinDoc.Find(".skinwrapper").Each(func(i int, selection *goquery.Selection) {
		img, _ := url.QueryUnescape(selection.Find(".charimg").First().Nodes[0].FirstChild.Attr[1].Val)
		skinName := selection.Find(".charnameEn").Text()
		compileRegex := regexp.MustCompile("_(.*?)_")
		match := compileRegex.FindStringSubmatch(img)
		if len(match) > 1 {
			name := match[1]
			skinCount[name] = append(skinCount[name], skinName)
		}
	})

	for i, operator := range operators {
		name := operator.Name
		if name == "阿米娅" {
			// 立绘
			for e := 0; e < 2; e++ {
				paintingName := fmt.Sprintf("立绘_%s_%d.png", name, e+1)
				m := hashutil.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
				var skin model.Skin
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
			// 精1立绘
			paintingName := fmt.Sprintf("立绘_%s_1+.png", name)
			m := hashutil.Md5(paintingName)
			path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
			painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
			var skin model.Skin
			skin.Url = painting
			operators[i].Skins = append(operators[i].Skins, skin)
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := hashutil.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
				var skin model.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		} else if name == "阿米娅(近卫)" || name == "阿米娅(医疗)" {
			// 立绘
			paintingName := fmt.Sprintf("立绘_%s_2.png", name)
			m := hashutil.Md5(paintingName)
			path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
			painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
			var skin model.Skin
			skin.Url = painting
			operators[i].Skins = append(operators[i].Skins, skin)
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := hashutil.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
				var skin model.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		}
		if operator.Rarity < 3 {
			// 立绘
			paintingName := fmt.Sprintf("立绘_%s_1.png", name)
			m := hashutil.Md5(paintingName)
			path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
			painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
			var skin model.Skin
			skin.Url = painting
			operators[i].Skins = append(operators[i].Skins, skin)
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := hashutil.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
				var skin model.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		} else if operator.Rarity >= 3 && !strings.Contains(operator.Name, "阿米娅") {
			// 立绘
			for e := 0; e < 2; e++ {
				paintingName := fmt.Sprintf("立绘_%s_%d.png", name, e+1)
				m := hashutil.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
				var skin model.Skin
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
			// 皮肤
			for c, sk := range skinCount[name] {
				paintingName := fmt.Sprintf("立绘_%s_skin%d.png", name, len(skinCount[name])-c)
				m := hashutil.Md5(paintingName)
				path := "https://media.prts.wiki" + fmt.Sprintf("/%s/%s/", m[:1], m[:2])
				painting := path + localassets.EscapePath(paintingName) + "?image_process=format,webp/quality,Q_90"
				var skin model.Skin
				skin.Name = sk
				skin.Url = painting
				operators[i].Skins = append(operators[i].Skins, skin)
			}
		}
		// 生日：新干员，以及生日还不是日期的干员（未公开、未录入、未知等）都重新查一次，
		// 这样干员页后续补上生日时缓存也能跟着更新
		config.DataMu.RLock()
		_, ignoreBirthday := config.IgnoreBirthday[name]
		config.DataMu.RUnlock()
		bucket, hasBirthday := birthdayBuckets[name]
		if (!hasBirthday || !isBirthdayDate(bucket)) && !ignoreBirthday {
			if birthday, ok := fetchBirthday(api, name); ok {
				removeFromBirthday(birthdays, bucket, name)
				birthdays[birthday] = append(birthdays[birthday], operators[i])
				birthdayBuckets[name] = birthday
			}
		}
	}

	// 生日缓存写回：日期变化的干员已经挪到新分组，空掉的分组直接删掉
	for bucket, list := range birthdays {
		key := birthdayKeyPrefix + bucket
		if len(list) == 0 {
			cache.RedisDel(key)
			continue
		}
		cache.RedisSet(key, json.MustMarshalString(list), 0)
	}

	cache.RedisSet("operatorList", json.MustMarshalString(operators), 0)
	log.Println("数据源更新完毕")
	// 开启本地素材缓存时，把干员的头像、半身像、立绘、皮肤增量下载到本地磁盘，
	// 后台执行不阻塞数据源更新，避免渲染图片时因为图床抖动或网络卡顿丢图
	if localassets.Enabled() && localassets.SyncAsync(localAssetURLs(operators)) {
		log.Println("已开始同步本地素材")
	}
	search.SetDataNeedUpdate()
	return nil
}

// localAssetURLs 收集需要缓存到本地的素材地址：
// 干员头像、半身像、立绘/皮肤（图床），敌人头像，以及游戏内皮肤素材（基建/干员箱页面用的头像与半身像）
func localAssetURLs(operators []model.Operator) []string {
	urls := make([]string, 0, len(operators)*6)
	for _, operator := range operators {
		urls = append(urls, operator.Avatar, operator.ThumbURL)
		for _, skin := range operator.Skins {
			urls = append(urls, skin.Url)
		}
	}
	urls = append(urls, search.EnemyAvatarURLs()...)
	urls = append(urls, skinAssetURLs()...)
	return urls
}

// skinAssetURLs 收集游戏内皮肤素材地址。
// 这些地址由皮肤表里的 skinId 拼出来（和模板里 {{urlquery .SkinId}} 的拼法一致），
// 皮肤表拿不到时返回空，不影响其它素材的同步。
func skinAssetURLs() []string {
	api := viper.GetString("api.skin_table")
	if api == "" {
		return nil
	}
	response, err := httpx.Open(api)
	if err != nil {
		log.Println("获取皮肤数据失败:", err)
		return nil
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println("读取皮肤数据失败:", err)
		return nil
	}
	var urls []string
	gjson.ParseBytes(body).Get("charSkins").ForEach(func(key, value gjson.Result) bool {
		urls = append(urls, skinAssetURLsOf(key.String())...)
		return true
	})
	return urls
}

// skinAssetURLsOf 生成单个皮肤的游戏内素材地址。
// 装置（trap_）和召唤物（token_）不会出现在任何渲染场景里，图床也没有对应素材，直接排除，
// 免得每次 /update 都去下上千个用不到、还会 404 的地址；真要用到时由 /local-assets 路由按需补齐。
func skinAssetURLsOf(skinId string) []string {
	if skinId == "" || strings.HasPrefix(skinId, "trap_") || strings.HasPrefix(skinId, "token_") {
		return nil
	}
	escaped := url.QueryEscape(skinId)
	return []string{
		"https://web.hycdn.cn/arknights/game/assets/char_skin/avatar/" + escaped + ".png",
		"https://web.hycdn.cn/arknights/game/assets/char_skin/portrait/" + escaped + ".png",
	}
}

// minOperators 干员一览解析结果的最小数量，低于这个数说明页面结构变了或拿到的是错误页
const minOperators = 50

// fetchDocument 抓取并解析页面：请求出错、状态码异常或解析失败时返回 false。
// 以前这里忽略 error 直接在 response.Body 上取内容，网络不通时 response 为 nil，
// 就会以「nil pointer dereference」的形式把整个 /update 崩掉。
func fetchDocument(url string) (*goquery.Document, bool) {
	response, err := httpx.Open(url)
	if err != nil {
		log.Println("获取页面失败:", url, err)
		return nil, false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		log.Printf("获取页面失败，状态码：%d，地址：%s", response.StatusCode, url)
		return nil, false
	}
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Println("解析页面失败:", url, err)
		return nil, false
	}
	return doc, true
}
