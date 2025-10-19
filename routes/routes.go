package routes

import (
	adminController "github.com/ahmadalaik/desa-digital/controllers/admin"
	authController "github.com/ahmadalaik/desa-digital/controllers/auth"
	"github.com/ahmadalaik/desa-digital/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/api")
	auth.POST("/login", authController.Login)

	// require authentication
	protected := router.Group("/api/admin")
	protected.Use(middlewares.AuthMiddleware())
	protected.GET("/dashboard", middlewares.Permission("dashboard-index"), adminController.Dashboard)

	return router
}
