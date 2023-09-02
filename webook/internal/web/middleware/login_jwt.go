package middleware

import (
	"awesomeProject/webook/internal/web"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"strings"
	"time"
)

//LoginMiddlewareBuilder扩展性

type LoginJWTMiddlewareBuilder struct {
	paths []string
}

func NewLoginJWTMiddlewareBuilder() *LoginJWTMiddlewareBuilder {
	return &LoginJWTMiddlewareBuilder{}
}
func (l *LoginJWTMiddlewareBuilder) IgnorePaths(path string) *LoginJWTMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginJWTMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//不需要登录校验
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		//sess := sessions.Default(ctx)
		//id := sess.Get("userID")
		//if id == nil {
		//	//没有登录
		//	ctx.AbortWithStatus(http.StatusUnauthorized)
		//	return
		//}
		//jwt实现
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 {
			//有人修改了Authorization请求头的值
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr := segs[1]
		claims := &web.UserCliams{}
		//ParseWithCliams里面，一定要传入指针
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		})

		//token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		//	return []byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"), nil
		//})
		if err != nil {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if token == nil || !token.Valid || claims.Uid == 0 {
			//没登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if claims.UserAgent != ctx.Request.UserAgent() {
			//ctx.String(http.StatusUnauthorized, "ua不对")
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		now := time.Now()
		if claims.ExpiresAt.Sub(now) < time.Second*50 {
			claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute))
			tokenStr, err = token.SignedString([]byte("95osj3fUD7fo0mlYdDbncXz4VD2igvf0"))
			if err != nil {
				//fmt.Printf("jwt失败", err)
				log.Println("jwt续约失败", err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}
		ctx.Set("claims", claims)
	}

}

//var IgnorePaths []string
//
//func CheckLogin()gin.HandlerFunc  {
//	return func(ctx *gin.Context) {
//
//	}
//}
