package seeders

import (
	"github.com/ahmadalaik/desa-digital/models"
	"gorm.io/gorm"
)

func SeedRoles(db *gorm.DB) {
	roles := []models.Role{
		{Name: "admin"},
		{Name: "user"},
	}

	for _, role := range roles {
		db.FirstOrCreate(&role, models.Role{Name: role.Name})

		var allPermissions []models.Permission
		db.Find(&allPermissions)

		switch role.Name {
		case "admin":
			db.Model(&role).Association("Permissions").Replace(allPermissions)
		case "user":
			var viewOnly []models.Permission
			db.Where("name IN ?", []string{"posts-index", "photos-index", "sliders-index", "pages-index"}).Find(&viewOnly)
			db.Model(&role).Association("Permissions").Replace(viewOnly)
		}

	}
}
