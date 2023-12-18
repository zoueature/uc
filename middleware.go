package uc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserAuthGinMiddleware gin 用户过滤中间件
func UserAuthGinMiddleware(decoder JwtEncoder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		forbidden := func() {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "Access Denied"})
		}
		token := getToken(ctx)
		if token == "" {
			forbidden()
			return
		}
		user, err := decoder.decodeJwt(token)
		if err != nil {
			forbidden()
			return
		}
		ctx.Set("user_id", user.Id)
		ctx.Set("user", user)
	}
}

// UserShouldAuthGinMiddleware 用户登录信息解析， 不强制登录
func UserShouldAuthGinMiddleware(decoder JwtEncoder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token := getToken(ctx)
		if token != "" {
			user, err := decoder.decodeJwt(token)
			if err == nil {
				ctx.Set("user_id", user.Id)
				ctx.Set("user", user)
			}
		}
	}
}

func getToken(ctx *gin.Context) string {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		token, _ = ctx.GetQuery("token")
	}
	if token == "" {
		token, _ = ctx.GetPostForm("token")
	}
	if token == "" {
		token, _ = ctx.Cookie("token")
	}
	return token
}
