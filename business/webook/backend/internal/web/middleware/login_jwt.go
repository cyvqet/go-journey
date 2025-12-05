package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginJwtMiddlewareBuilder struct {
	paths []string
}

func NewLoginJwtMiddlewareBuilder() *LoginJwtMiddlewareBuilder {
	return &LoginJwtMiddlewareBuilder{}
}

func (l *LoginJwtMiddlewareBuilder) IgnorePath(paths string) *LoginJwtMiddlewareBuilder {
	l.paths = append(l.paths, paths)
	return l
}

func (l *LoginJwtMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if slices.Contains(l.paths, ctx.Request.RequestURI) {
			ctx.Next()
			return
		}

		// JWT 验证逻辑
		// 从请求头中获取 JWT tokenHeader，进行验证
		// 如果验证失败，返回 401 未授权错误
		// 如果验证成功，调用 ctx.Next() 继续处理请求
		tokenHeader := ctx.GetHeader("Authorization")
		if tokenHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			return
		}

		segs := strings.Split(tokenHeader, " ")
		if len(segs) != 2 || segs[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			return
		}

		token := segs[1]
		jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil
		})

		if err != nil || !jwtToken.Valid {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			return
		}

		fmt.Printf("JWT 验证通过: %+v\n", jwtToken.Claims)

		ctx.Next()
	}
}
