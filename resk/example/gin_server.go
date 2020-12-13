package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	app := gin.Default()
	//app.GET("/hello", func(context *gin.Context) {
	//	context.String(http.StatusOK, "hello gin")
	//})
	//app.GET("/user/:name", func(context *gin.Context) {
	//	name := context.Param("name")
	//	context.String(http.StatusOK, name)
	//})
	app.Run(":8086")
}
