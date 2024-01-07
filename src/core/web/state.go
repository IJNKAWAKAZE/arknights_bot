package web

import (
	"arknights_bot/plugins/skland"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"os"
)

func State(r *gin.Engine) {
	r.GET("/state/:data/:uid", func(c *gin.Context) {
		/*var account skland.Account
		uid := c.Param("uid")
		json.Unmarshal([]byte(c.Param("data")), &account)
		playerData, err := skland.GetPlayerInfo(uid, account)
		if err != nil {
			log.Println(err)
			return
		}*/
		open, err := os.Open("player_info.txt")
		if err != nil {
			return
		}
		var aaa skland.PlayerData
		readAll, _ := io.ReadAll(open)
		playerData := string(readAll)
		json.Unmarshal([]byte(playerData), &aaa)
		c.HTML(http.StatusOK, "State.tmpl", aaa)
	})
}
