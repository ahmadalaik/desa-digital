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
	// param1 url, param2 middleware checks whether the user has permission, param3 function (controller)
	protected.GET("/dashboard", middlewares.Permission("dashboard-index"), adminController.Dashboard)

	protected.GET("/permissions", middlewares.Permission("permissions-index"), adminController.FindPermissons)
	protected.POST("/permissions", middlewares.Permission("permissions-create"), adminController.CreatePermission)
	protected.GET("/permissions/:id", middlewares.Permission("permissions-show"), adminController.FindPermissonByID)
	protected.PUT("/permissions/:id", middlewares.Permission("permissions-update"), adminController.UpdatePermission)
	protected.DELETE("/permissions/:id", middlewares.Permission("permissions-delete"), adminController.DeletePermission)
	protected.GET("/permissions/all", middlewares.Permission("permissions-index"), adminController.FindAllPermissions)

	// role routes
	protected.GET("/roles", middlewares.Permission("roles-index"), adminController.FindRoles)
	protected.POST("/roles", middlewares.Permission("roles-create"), adminController.CreateRole)
	protected.GET("/roles/:id", middlewares.Permission("roles-show"), adminController.FindRoleByID)
	protected.PUT("/roles/:id", middlewares.Permission("roles-update"), adminController.UpdateRole)
	protected.DELETE("/roles/:id", middlewares.Permission("roles-delete"), adminController.DeleteRole)
	protected.GET("/roles/all", middlewares.Permission("roles-index"), adminController.FindAllRoles)

	// user routes
	protected.GET("/users", middlewares.Permission("users-index"), adminController.FindUsers)
	protected.POST("/users", middlewares.Permission("users-create"), adminController.CreateUser)
	protected.GET("/users/:id", middlewares.Permission("users-show"), adminController.FindUserByID)
	protected.PUT("/users/:id", middlewares.Permission("users-Update"), adminController.UpdateUser)
	protected.DELETE("/users/:id", middlewares.Permission("users-delete"), adminController.DeleteUser)

	return router
}
