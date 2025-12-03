package main

import (
	"log"
	"time"
	"webook/internal/repository"
	"webook/internal/repository/dao"
	"webook/internal/service"
	"webook/internal/web"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := initWebServer()
	db := initDB()
	userHandler := initUser(db)

	userHandler.RegisterRouter(server)

	err := server.Run(":8080")
	if err != nil {
		log.Println("Web 服务启动失败:", err)
		return
	}
}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}

	return db
}

func initWebServer() *gin.Engine {
	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://foo.com"}, // 指定允许跨域的域名列表

		AllowMethods: []string{"GET", "POST"}, // 指定允许跨域的 HTTP 方法

		AllowHeaders: []string{"Origin", "Authorization", "Content-Type"}, // 指定浏览器可以携带的请求头

		ExposeHeaders: []string{"Content-Length"}, // 指定哪些响应头可以暴露给前端 JS

		// AllowCredentials: 是否允许跨域携带 cookies 或者 Authorization 等凭证
		// 注意：开启后，AllowOrigins 不能为 "*"（必须明确指定域名）
		AllowCredentials: true,

		// AllowOriginFunc: 自定义判断逻辑（动态允许的 Origin）
		// 这里表示：如果 Origin == "https://github.com"，也允许
		// 优先级：这个函数的判断会比 AllowOrigins 更高
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},

		MaxAge: 12 * time.Hour, // 浏览器对预检请求（OPTIONS）的缓存时间
	}))

	return server
}

func initUser(db *gorm.DB) *web.UserHandler {
	userDao := dao.NewUserDao(db)
	userRepository := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	return userHandler
}
