package seeders

import (
	"log"

	"github.com/ahmadalaik/desa-digital/database"
)

func Seed() {
	db := database.DB
	log.Println("Running database seeders...")

	SeedPermissions(db)
	SeedRoles(db)
	SeedUsers(db)

	log.Println("Database seeding completed!")
}
