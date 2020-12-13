package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	app := gin.Default()
	app.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "hello gin")
	})
	app.Run()
}
