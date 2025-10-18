package routes

import (
	"github.com/gin-gonic/gin"
	authController "github.com/ahmadalaik/desa-digital/controllers/auth"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	auth := router.Group("/api")
	auth.POST("/login", authController.Login)
	
	return router
}
