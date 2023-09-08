package ioc

import (
	"awesomeProject/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"strings"
	"time"
)

func InitWebserver(mdls []gin.HandlerFunc, hdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	hdl.RegisterRoutes(server)
	return server
}
func InitMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		cors.New(cors.Config{
			//AllowOrigins: []string{"*"},
			//AllowMethods: []string{"POST", "GET"},
			AllowHeaders: []string{"Content-Type", "Authorization"},
			// 你不加这个，前端是拿不到的
			ExposeHeaders: []string{"x-jwt-token"},
			//ExposeHeaders: []string{"x-jwt-token"},
			// 是否允许你带 cookie 之类的东西
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				if strings.HasPrefix(origin, "http://localhost") {
					// 你的开发环境
					return true
				}
				return strings.Contains(origin, "yourcompany.com")
			},
			MaxAge: 12 * time.Hour,
		})}

}
