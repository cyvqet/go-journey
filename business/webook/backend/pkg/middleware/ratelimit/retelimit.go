package ratelimit

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Builder 限流器构造器，用于创建基于滑动窗口算法的限流中间件
type Builder struct {
	prefix   string        // Redis键前缀
	cmd      redis.Cmdable // Redis客户端
	interval time.Duration // 时间窗口长度
	rate     int           // 窗口内允许的最大请求数
}

//go:embed slide_window.lua
var luaScript string // 嵌入的滑动窗口Lua脚本

// NewBuilder 创建Builder实例
func NewBuilder(cmd redis.Cmdable, interval time.Duration, rate int) *Builder {
	return &Builder{
		cmd:      cmd,
		prefix:   "ip-limiter", // 默认前缀
		interval: interval,
		rate:     rate,
	}
}

// Prefix 设置Redis键前缀
func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

// Build 创建Gin限流中间件
func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			// Redis故障时的处理策略：保守做法（限流）vs 激进做法（放行）
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

// limit 执行限流检查
func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	// 基于客户端IP构造Redis键
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	// 执行Lua脚本实现原子化的滑动窗口限流
	return b.cmd.Eval(ctx, luaScript, []string{key},
		b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
