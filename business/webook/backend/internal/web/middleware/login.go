package middleware

import (
	"log"
	"net/http"
	"slices"
	"time"

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

		session.Options(sessions.Options{
			MaxAge: 3600,
		})

		// 刷新会话 session 的过期时间
		updateTime := session.Get("update_time")
		now := time.Now().UnixMilli()
		if updateTime == nil {
			log.Println("第一次刷新会话时间")
			session.Set("update_time", now)
			session.Save()
			ctx.Next()
			return
		}

		updateTimeValue, ok := updateTime.(int64)
		if !ok {
			log.Println("会话时间格式错误")
			ctx.JSON(http.StatusInternalServerError, gin.H{"message": "系统错误"})
			ctx.Abort()
			return
		}

		if now-updateTimeValue > 60*1000 { // 1分钟没有操作，刷新 session
			log.Println("刷新会话时间")
			session.Set("update_time", now)
			session.Save()
		}

		log.Printf("当前会话的用户邮箱: %v\n", email)
		ctx.Next()
	}
}
