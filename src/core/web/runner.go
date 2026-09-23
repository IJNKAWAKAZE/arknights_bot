package web

import (
	"arknights_bot/utils/localassets"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"sync"
)

var httpServer *http.Server
var serverMu sync.Mutex
var stopping bool

func Start() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	// 本地素材缓存：开启后页面里的图床地址会改写成同源的 /local-assets 地址
	localassets.Register(r)
	r.Static("/assets", "./assets")
	r.Static("/template/js", "./template/js")
	Help(r)
	State(r)
	Box(r)
	Missing(r)
	Gacha(r)
	Card(r)
	Base(r)
	Headhunt(r)
	Operator(r)
	Enemy(r)
	Recruit(r)
	Calendar(r)
	Depot(r)
	BoxDetail(r)
	Summary(r)
	Lottery(r)
	host := viper.GetString("http.host")
	if host == "" {
		host = "127.0.0.1"
	}
	addr := host + ":" + viper.GetString("http.port")
	serverMu.Lock()
	if stopping {
		serverMu.Unlock()
		return
	}
	httpServer = &http.Server{Addr: addr, Handler: r}
	server := httpServer
	serverMu.Unlock()
	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// Shutdown 优雅关闭 Web 服务
func Shutdown(ctx context.Context) error {
	serverMu.Lock()
	stopping = true
	server := httpServer
	serverMu.Unlock()
	if server == nil {
		return nil
	}
	if err := server.Shutdown(ctx); err != nil {
		log.Println("Web服务关闭失败:", err)
		return err
	}
	return nil
}
