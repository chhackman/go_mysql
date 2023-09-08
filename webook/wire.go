//go:build wireinject

package main

import (
	"awesomeProject/webook/internal/repository"
	cache2 "awesomeProject/webook/internal/repository/cache"
	"awesomeProject/webook/internal/repository/dao"
	"awesomeProject/webook/internal/service"
	"awesomeProject/webook/internal/web"
	"awesomeProject/webook/ioc"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	//"github.com/google/wire"
)

func InitWebServer() *gin.Engine {

	wire.Build(
		//最基础的第三方依赖
		ioc.InitDB, ioc.InitRedis,
		//初始化DAO
		dao.NewUserDAO,

		cache2.NewUserCache,
		cache2.NewCodeCache,

		repository.NewCodeRepository,
		repository.NewUserRepository,

		service.NewCodeService,
		service.NewUserService,

		//基于内存的实现
		ioc.InitSMSService,

		web.NewUserHandler,

		//gin.Default,

		ioc.InitWebserver,
		ioc.InitMiddlewares,
	)
	return gin.Default()
}
