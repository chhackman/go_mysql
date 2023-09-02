package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"Blog": "www.flysnow.org",
			"Name": "wangjie",
		})
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		c.String(http.StatusOK, "the user id is %s", id)
	})

	//r.GET("/users/*id", func(c *gin.Context) {
	//	id := c.Param("id")
	//	c.String(http.StatusOK, "the user id is %s", id)
	//})

	r.Run(":8080")
}
