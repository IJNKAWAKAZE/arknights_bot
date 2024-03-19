package web

import (
	"arknights_bot/plugins/material"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Material(r *gin.Engine) {
	r.GET("/material", func(c *gin.Context) {
		r.LoadHTMLFiles("./template/Material.tmpl")
		name := c.Query("name")
		materialInfo := material.ParseMaterial(name)
		c.HTML(http.StatusOK, "Material.tmpl", materialInfo)
	})
}
