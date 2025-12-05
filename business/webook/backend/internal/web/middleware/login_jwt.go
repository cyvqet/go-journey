package middleware

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"
	"webook/internal/web"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// token有效期：30分钟，剩余5分钟时刷新
const (
	// tokenTTL    = 30 * time.Minute
	refreshWhen = 5 * time.Minute  // 剩余5分钟刷新
	newTokenTTL = 30 * time.Minute // 新token还是30分钟
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
		claim := web.UserClaims{} // 自定义的 Claims 结构体
		token := segs[1]
		jwtToken, err := jwt.ParseWithClaims(token, &claim, func(t *jwt.Token) (any, error) {
			return []byte("secret"), nil
		})

		if err != nil || !jwtToken.Valid { // token 过期 Valid 返回 false
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "未授权"})
			return
		}

		remaining := time.Until(claim.ExpiresAt.Time)
		fmt.Printf("token 剩余时间: %v\n", remaining)
		if remaining <= refreshWhen {
			fmt.Printf("token 即将过期，剩余时间: %v，刷新 token\n", remaining)
			// 刷新 token 的过期时间
			claim.ExpiresAt = jwt.NewNumericDate(time.Now().Add(newTokenTTL))

			// 生成新的 JWT token
			newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
			newTokenStr, err := newToken.SignedString([]byte("secret"))
			if err != nil {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
				return
			}

			// 将新的 token 返回给客户端
			ctx.Header("Jwt-Token", newTokenStr)
		}

		// 将 claim 存储到上下文中，供后续处理使用
		ctx.Set("claim", claim)

		ctx.Next()
	}
}
