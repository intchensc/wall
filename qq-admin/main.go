package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/**/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login/login.html", gin.H{
			"title": "登录测试页面",
			"data":  "good circle",
		})
	})
	r.Run(":9090")
}
