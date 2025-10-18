package main

import (
	"github.com/ahmadalaik/desa-digital/config"
	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/database/seeders"
	"github.com/ahmadalaik/desa-digital/routes"
)

func main() {
	config.LoadEnv()

	database.InitDB()

	seeders.Seed()

	r := routes.SetupRouter()

	r.Run(":" + config.GetEnv("APP_PORT", "8080"))
}
