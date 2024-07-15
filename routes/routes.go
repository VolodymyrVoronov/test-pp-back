package routes

import (
	"test-pp-back/handlers"
	"test-pp-back/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {

	server.Use(middlewares.EnableCORS())

	auth := server.Group("/auth")
	auth.POST("/register", handlers.RegisterHandler)
	auth.POST("/login", handlers.LoginHandler)
	auth.GET("/verify", handlers.VerifyHandler)
	auth.GET("/check-auth", handlers.CheckAuthHandler)

}
