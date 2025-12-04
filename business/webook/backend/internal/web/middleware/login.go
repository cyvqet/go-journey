package middleware

import (
	"log"
	"net/http"
	"slices"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePath(paths string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(l.paths, ctx.Request.RequestURI) {
			ctx.Next()
			return
		}

		session := sessions.Default(ctx)
		email := session.Get("userEmail")
		if email == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"message": "未登录"})
			ctx.Abort()
			return
		}
		log.Printf("当前会话的用户邮箱: %v\n", email)
		ctx.Next()
	}
}
