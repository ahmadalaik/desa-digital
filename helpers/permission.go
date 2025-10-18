package helpers

import "github.com/ahmadalaik/desa-digital/models"

func GetPermission(roles []models.Role) map[string]bool {
	permissionMap := make(map[string]bool)

	for _, role := range roles {
		for _, permission := range role.Permissions {
			permissionMap[permission.Name] = true
		}
	}
	return permissionMap
}
