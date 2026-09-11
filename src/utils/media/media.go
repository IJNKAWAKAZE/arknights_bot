package media

import (
	"bytes"
	"fmt"
	"github.com/mxschmitt/playwright-go"
	"github.com/spf13/viper"
	"github.com/tidwall/gjson"
	"golang.org/x/image/webp"
	"image"
	"image/color"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"sync"
	"time"
)

var (
	browser     playwright.Browser
	contexts    = make(map[float64]playwright.BrowserContext)
	browserLock sync.Mutex
)

// browserContext 返回可复用的浏览器上下文。
// 复用同一个上下文可以让远程图片命中浏览器缓存，避免每次截图都重新下载图床资源，
// 从而减少加载超时和占位图。
func browserContext(scale float64) (playwright.BrowserContext, error) {
	browserLock.Lock()
	defer browserLock.Unlock()
	if browser != nil && !browser.IsConnected() {
		browser = nil
		contexts = make(map[float64]playwright.BrowserContext)
	}
	if browser == nil {
		pw, err := playwright.Run()
		if err != nil {
			log.Println("未检测到playwright，开始自动安装...")
			if installErr := playwright.Install(&playwright.RunOptions{Browsers: []string{"chromium"}}); installErr != nil {
				return nil, fmt.Errorf("playwright安装失败: %w", installErr)
			}
			pw, err = playwright.Run()
			if err != nil {
				return nil, fmt.Errorf("playwright启动失败: %w", err)
			}
		}
		browser, err = pw.Chromium.Launch()
		if err != nil {
			browser = nil
			log.Println(err)
			return nil, fmt.Errorf("playwright启动失败: %w", err)
		}
	}
	if ctx, ok := contexts[scale]; ok && ctx != nil && !ctx.IsClosed() {
		return ctx, nil
	}
	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{DeviceScaleFactor: &scale})
	if err != nil {
		return nil, fmt.Errorf("创建浏览器上下文失败: %w", err)
	}
	contexts[scale] = ctx
	return ctx, nil
}

// waitForResources 等待页面上的字体与图片加载完成。
// 加载失败的图片会先重试一次，只有确认失败才用占位图兜底，并锁定渲染尺寸避免撑开布局。
const waitForResources = `async (timeoutMs) => {
    const FALLBACK = new URL("/assets/common/amiya.png", window.location.href).href;
    const sleep = (ms) => new Promise((resolve) => setTimeout(resolve, ms));
    const isLoaded = (img) => img.complete && img.naturalWidth > 0;
    const isFailed = (img) => img.complete && img.naturalWidth === 0;

    const images = Array.from(document.images).filter((img) => !!(img.getAttribute("src") || img.src));
    if (images.length === 0) {
        return true;
    }

    // 接管错误处理，避免模板里的 onerror 在图片稍慢时就立刻替换成占位图。
    for (const img of images) {
        img.onerror = null;
        img.removeAttribute("onerror");
    }

    const retried = new WeakSet();
    const settled = (img) => isLoaded(img) || img.src === FALLBACK || (isFailed(img) && retried.has(img));

    let fontsReady = !(document.fonts && document.fonts.ready);
    if (!fontsReady) {
        document.fonts.ready.then(() => { fontsReady = true; }, () => { fontsReady = true; });
    }

    // 等待所有图片真正加载完成；失败的先重试一次，只有确认失败才用占位图兜底。
    const deadline = Date.now() + timeoutMs;
    while (Date.now() < deadline) {
        for (const img of images) {
            if (isFailed(img) && !retried.has(img) && img.src !== FALLBACK) {
                retried.add(img);
                const src = img.src || img.getAttribute("src");
                if (src) {
                    img.src = src + (src.indexOf("?") === -1 ? "?" : "&") + "_retry=" + Date.now();
                }
            }
        }
        if (fontsReady && images.every(settled)) {
            break;
        }
        await sleep(150);
    }

    // 仍然不可用的图片用占位图兜底，并锁定当前渲染尺寸，避免大占位图撑开布局。
    const broken = images.filter((img) => !isLoaded(img));
    for (const img of broken) {
        const rect = img.getBoundingClientRect();
        if (rect.width > 0) { img.style.width = Math.round(rect.width) + "px"; }
        if (rect.height > 0) { img.style.height = Math.round(rect.height) + "px"; }
        if (rect.width <= 0 && rect.height <= 0) {
            // 无法从布局中获得尺寸时，限制占位图大小，避免 512x512 的占位图撑开页面。
            img.style.maxWidth = "120px";
            img.style.maxHeight = "120px";
        }
        img.style.objectFit = "cover";
        if (img.src !== FALLBACK) { img.src = FALLBACK; }
    }
    if (broken.length > 0) {
        const fallbackDeadline = Date.now() + 3000;
        while (Date.now() < fallbackDeadline && !broken.every(isLoaded)) {
            await sleep(150);
        }
    }

    await Promise.all(images.map(async (img) => {
        if (isLoaded(img) && img.decode) {
            try { await img.decode(); } catch (_) {}
        }
    }));
    return true;
}`

// imageTimeoutMs 返回截图时等待远程图片加载的最长时间（毫秒）。
// 数值越大越不容易出现占位图，但会占用截图任务更久，因此做成可配置项：
// screenshot.image_timeout（秒），默认 10 秒。
func imageTimeoutMs() int {
	seconds := viper.GetInt("screenshot.image_timeout")
	if seconds <= 0 {
		seconds = 10
	}
	return seconds * 1000
}

// Screenshot 屏幕截图
func Screenshot(url string, waitTime float64, scale float64) ([]byte, error) {
	context, err := browserContext(scale)
	if err != nil {
		return nil, err
	}
	page, err := context.NewPage()
	if err != nil {
		return nil, fmt.Errorf("创建页面失败: %w", err)
	}
	defer func() {
		log.Println("关闭playwright")
		page.Close()
	}()
	log.Println("开始进行截图...")
	resp, err := page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateDomcontentloaded,
	})
	if err != nil {
		return nil, fmt.Errorf("页面加载失败: %w", err)
	}
	if resp != nil && resp.Status() >= 400 {
		return nil, fmt.Errorf("页面加载失败，状态码：%d", resp.Status())
	}
	// 等待字体与图片加载完成，加载失败的图片重试后再用占位图兜底。
	// 等待时长可配置：值越大越完整，代价是这张截图占用更久。
	timeoutMs := imageTimeoutMs()
	if _, err := page.WaitForFunction(waitForResources, timeoutMs, playwright.PageWaitForFunctionOptions{Timeout: playwright.Float(float64(timeoutMs + 8000))}); err != nil {
		log.Println("等待图片和字体加载超时，继续截图:", err)
	}
	page.WaitForTimeout(waitTime)
	locator := page.Locator("#main")
	if v, err := locator.IsVisible(); err != nil || !v {
		log.Println("元素未加载取消截图操作")
		return nil, fmt.Errorf("元素未加载")
	}
	// 检测错误页，把错误文本作为截图失败原因返回
	if class, _ := locator.GetAttribute("class"); strings.Contains(class, "error") {
		if text, err := locator.InnerText(); err == nil && text != "" {
			return nil, fmt.Errorf("%s", text)
		}
	}
	screenshot, err := locator.Screenshot(playwright.LocatorScreenshotOptions{Type: playwright.ScreenshotTypeJpeg})
	if err != nil {
		return nil, fmt.Errorf("截图失败: %w", err)
	}
	log.Println("截图完成...")
	return screenshot, nil
}

var imgClient = &http.Client{Timeout: 15 * time.Second}

func GetImg(url string) []byte {
	var pic []byte
	times := 0
	for times < 3 {
		resp, err := imgClient.Get(url)
		if err != nil {
			log.Println("获取图片失败", err)
			times++
			continue
		}
		if resp.StatusCode != http.StatusOK {
			log.Printf("获取图片失败，状态码：%d，URL：%s", resp.StatusCode, url)
			resp.Body.Close()
			times++
			continue
		}
		pic, err = io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			log.Println("读取图片失败", err)
			times++
			continue
		}
		break
	}
	return pic
}

func ImgConvert(url string) []byte {
	pic, err := http.Get(url)
	if err != nil {
		log.Println("获取图片失败", err)
		return nil
	}
	m, err := webp.Decode(pic.Body)
	if err != nil {
		log.Println("解析图片失败", err)
		return nil
	}
	bounds := m.Bounds()
	dx := bounds.Dx()
	dy := bounds.Dy()
	newRgba := image.NewRGBA(bounds)
	f := true
	go overtime(&f)
o:
	for i := 0; i < dx; i++ {
		for j := 0; j < dy; j++ {
			if !f {
				log.Println("图片转换超时")
				break o
			}
			colorRgb := m.At(i, j)
			r, g, b, a := colorRgb.RGBA()
			r_uint8 := uint8(r >> 8)
			g_uint8 := uint8(g >> 8)
			b_uint8 := uint8(b >> 8)
			a_uint8 := uint8(a >> 8)

			if a_uint8 > 23 {
				r_uint8 = 255
				g_uint8 = 255
				b_uint8 = 255
			}
			a_uint8 = 255
			newRgba.SetRGBA(i, j, color.RGBA{R: r_uint8, G: g_uint8, B: b_uint8, A: a_uint8})
		}
	}
	if !f {
		return nil
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, newRgba)
	return buf.Bytes()
}

// CutImg 图片裁剪
func CutImg(url string) []byte {
	pic, err := http.Get(url)
	if err != nil {
		log.Println("获取图片失败", err)
		return nil
	}

	var subImage image.Image
	if pic.Header.Get("Content-Type") == "image/webp" {
		m, err := webp.Decode(pic.Body)
		if err != nil {
			log.Println("解析图片失败", err)
			return nil
		}
		rgba := m.(*image.NYCbCrA)
		subImage = rgba.SubImage(image.Rect(0, m.Bounds().Dy(), m.Bounds().Dx(), int(float64(m.Bounds().Dy())/1.5))).(*image.NYCbCrA)
	} else {
		m, _, err := image.Decode(pic.Body)
		if err != nil {
			log.Println("解析图片失败", err)
			return nil
		}
		rgba := m.(*image.NRGBA)
		subImage = rgba.SubImage(image.Rect(0, m.Bounds().Dy(), m.Bounds().Dx(), int(float64(m.Bounds().Dy())/1.5))).(*image.NRGBA)
	}
	buf := new(bytes.Buffer)
	png.Encode(buf, subImage)
	return buf.Bytes()
}

func overtime(f *bool) {
	time.Sleep(time.Second * 10)
	*f = false
}

// OCR OCR识别
func OCR(file io.Reader, lang, engine, sep string) ([]string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", "test.png")
	if err != nil {
		log.Println("创建文件失败")
		return nil, err
	}
	io.Copy(part, file)
	writer.WriteField("language", lang)
	writer.WriteField("FileType", ".Auto")
	writer.WriteField("OCREngine", engine)
	writer.Close()
	request, err := http.NewRequest("POST", "https://api.ocr.space/parse/image", body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	request.Header.Add("Apikey", "helloworld")
	client := http.DefaultClient
	client.Timeout = time.Second * 10
	resp, err := client.Do(request)
	if err != nil {
		log.Println("ocr失败")
		return nil, err
	}
	read, _ := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	result := gjson.ParseBytes(read)
	log.Println("识别结果：", result.String())
	return strings.Split(result.Get("ParsedResults.0.ParsedText").String(), sep), nil
}

// CreateTelegraphPage 创建telegraph页面
func CreateTelegraphPage(content, title string) string {
	api := viper.GetString("api.telegraph")
	request, _ := http.NewRequest("GET", api, nil)
	params := request.URL.Query()
	params.Add("access_token", viper.GetString("telegraph.token"))
	params.Add("title", title)
	params.Add("content", content)
	request.URL.RawQuery = params.Encode()
	response, _ := http.DefaultClient.Do(request)
	readAll, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	jsonStr := string(readAll)
	log.Println(jsonStr)
	url := gjson.Get(jsonStr, "result.url").String()
	return url
}
