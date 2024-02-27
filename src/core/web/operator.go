package web

import (
	"arknights_bot/plugins/operator"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Operator(r *gin.Engine) {
	r.GET("/operator", func(c *gin.Context) {
		name := c.Query("name")
		operatorInfo := operator.ParseOperator(name)
		c.HTML(http.StatusOK, "Operator.tmpl", operatorInfo)
	})
}
