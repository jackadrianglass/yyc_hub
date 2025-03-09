package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.tmpl", gin.H{})
	})
	router.GET("/ping", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "ping.tmpl", gin.H{})
	
	})

	router.Run(":8080")
}
