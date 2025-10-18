package main

import (
	"github.com/ahmadalaik/desa-digital/config"
	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/database/seeders"
	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	database.InitDB()

	seeders.Seed()

	r := gin.Default()

	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Hello World",
		})
	})

	r.Run(":" + config.GetEnv("APP_PORT", "8080"))
}
