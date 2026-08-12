package web

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

var httpServer *http.Server

func Start() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
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
	httpServer = &http.Server{Addr: addr, Handler: r}
	err := httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}

// Shutdown 优雅关闭 Web 服务
func Shutdown() {
	if httpServer == nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Println("Web服务关闭失败:", err)
	}
}
