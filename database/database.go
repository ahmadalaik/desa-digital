package database

import (
	"fmt"
	"log"

	"github.com/ahmadalaik/desa-digital/config"
	"github.com/ahmadalaik/desa-digital/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	dbUser := config.GetEnv("DB_USER", "ahmadalaik")
	//dbPass := config.GetEnv("DB_PASS", "")
	dbHost := config.GetEnv("DB_HOST", "localhost")
	dbPort := config.GetEnv("DB_PORT", "5432")
	dbName := config.GetEnv("DB_NAME", "db_desa_digital")

	dsn := fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=disable", dbUser, dbHost, dbPort, dbName)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("failed connecting to database: ", err)
	}

	DB = db
	fmt.Println("Database connected successfully!")

	err = DB.AutoMigrate(&models.User{}, &models.Role{}, &models.Permission{}, &models.Category{}, &models.Post{}, &models.Slider{}, &models.Page{}, &models.Photo{}, &models.Aparatur{}, &models.Product{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	fmt.Println("Database migrate successfully!")
}
