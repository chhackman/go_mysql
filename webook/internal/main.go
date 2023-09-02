package main

import (
	"awesomeProject/webook/config"
	"awesomeProject/webook/internal/repository"
	cache2 "awesomeProject/webook/internal/repository/cache"
	"awesomeProject/webook/internal/repository/dao"
	"awesomeProject/webook/internal/service"
	"awesomeProject/webook/internal/web"
	"awesomeProject/webook/internal/web/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	redis2 "github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"strings"
	"time"
)

func main() {
	db := initDB()
	redisCon := redis2.NewClient(&redis2.Options{
		Addr:     config.Config.Redis.Addr,
		Password: "HzN8m%cr!Vve",
		DB:       12,
	})
	server := initWebServer()
	u := initUser(db, redisCon)
	u.RegisterRoutes(server)
	server.Run(":8081")
}

func initUser(db *gorm.DB, redisCmd *redis2.Client) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	cacheV3 := cache2.NewUserCache(redisCmd)
	repo := repository.NewUserRepository(ud, cacheV3)
	svc := service.NewUserService(repo)
	u := web.NewUserHandler(svc)
	return u
}
func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		println("这是第一个middleware")
	})
	server.Use(func(ctx *gin.Context) {
		println("这是第二个middleware")
	})

	server.Use(cors.New(cors.Config{
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
	}))

	////步骤1
	//store := cookie.NewStore([]byte("secret"))
	//server.Use(sessions.Sessions("mysession", store))
	//
	////步骤3
	//server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	//redis存储session
	print(config.Config.Redis.Addr)
	store, err := redis.NewStore(16, "tcp", config.Config.Redis.Addr, "HzN8m%cr!Vve", []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), []byte("0Pf2r0wZBpXVXlQNdpwCXN4ncnlnZSc3"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))

	server.Use(middleware.NewLoginMiddlewareBuilder().IgnorePaths("/users/signup").
		IgnorePaths("/users/login").Build())
	//jwt
	//server.Use(middleware.NewLoginJWTMiddlewareBuilder().IgnorePaths("/users/signup").
	//	IgnorePaths("/users/login").Build())
	return server
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("indigo:indigotest@tcp(10.1.80.122:3306)/go_test"))
	//db, err := gorm.Open(mysql.Open("indigo:indigotest@tcp(10.1.90.122:3306)/go_test"))
	//db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
