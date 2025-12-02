package main

import (
	"webook/internal/web"

	"github.com/gin-gonic/gin"
)

func main() {
	server := gin.Default()

	userHandler := web.NewUserHandler()
	userHandler.RegisterRouter(server)

	server.Run(":8080")
}
