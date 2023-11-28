package uc

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserAuthGinMiddleware gin 用户过滤中间件
func UserAuthGinMiddleware(decoder JwtEncoder) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		forbidden := func() {
			ctx.JSON(http.StatusForbidden, gin.H{"status": "Access Denied"})
		}
		token := ctx.GetHeader("Authorization")
		if token == "" {
			token, _ = ctx.GetQuery("token")
		}
		if token == "" {
			token, _ = ctx.GetPostForm("token")
		}
		if token == "" {
			forbidden()
			return
		}
		user, err := decoder.decodeJwt(token)
		if err != nil {
			forbidden()
			return
		}
		ctx.Set("userId", user.Id)
		ctx.Set("user", user)
	}
}
