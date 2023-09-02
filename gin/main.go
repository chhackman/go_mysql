package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.POST("/post", func(context *gin.Context) {
		context.String(http.StatusOK, "hello post")
	})

	r.GET("/user/:name", func(context *gin.Context) {
		name := context.Param("name")
		context.String(http.StatusOK, "hello,这是参数路由"+name)
	})

	r.GET("/views/*.html", func(context *gin.Context) {
		path := context.Param(".html")
		context.String(http.StatusOK, "这是通配符路由"+path)
	})

	r.GET("/order", func(context *gin.Context) {
		oid := context.Query("id")
		context.String(http.StatusOK, "这是查询参数"+oid)
	})

	//r.GET("/wangjie",)

	r.Run() // listen and serve on 0.0.0.0:8080
}
