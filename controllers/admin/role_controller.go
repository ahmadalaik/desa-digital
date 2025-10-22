package admin

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindRoles(c *gin.Context) {
	var roles []models.Role
	var roleResponse []structs.RoleResponse
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Preload("Permissions").Model(&models.Role{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&roles).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch roles",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	for _, role := range roles {
		permissionResponse := []structs.PermissionResponse{}

		for _, permission := range role.Permissions {
			permissionResponse = append(permissionResponse, structs.PermissionResponse{
				Id:        permission.ID,
				Name:      permission.Name,
				CreatedAt: permission.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: permission.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		roleResponse = append(roleResponse, structs.RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Permissions: permissionResponse,
			CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	helpers.PaginateResponse(c, roleResponse, total, page, limit, baseURL, search, "List Data Roles")
}

func CreateRole(c *gin.Context) {
	var req structs.RoleCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN ?", req.PermissionIDs).Find(&permissions)
	}

	role := models.Role{
		Name:        req.Name,
		Permissions: permissions,
	}

	if err := database.DB.Create(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create role",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Success create role",
		Data:    role,
	})
}

func FindRoleByID(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.Preload("Permissions").Where("id = ?", id).First(&role).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	permissionResponses := []structs.PermissionResponse{}
	for _, permission := range role.Permissions {
		permissionResponses = append(permissionResponses, structs.PermissionResponse{
			Id:        permission.ID,
			Name:      permission.Name,
			CreatedAt: permission.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: permission.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	roleResponse := structs.RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Permissions: permissionResponses,
		CreatedAt:   role.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Role found",
		Data:    roleResponse,
	})
}

func UpdateRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.Where("id = ?", id).First(&role).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.RoleUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	role.Name = req.Name

	var permissions []models.Permission
	if len(req.PermissionIDs) > 0 {
		database.DB.Where("id IN = ?", req.PermissionIDs).Find(&permissions)
	}
	database.DB.Model(&role).Association("Permissions").Replace(&permissions)

	if err := database.DB.Save(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update role",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success update role",
		Data:    role,
	})
}

func DeleteRole(c *gin.Context) {
	id := c.Param("id")
	var role models.Role

	if err := database.DB.First(&role, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Role not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Model(&role).Association("Permissions").Clear(); err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Failed to detach role from permissions",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&role).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete role",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete role",
		Data:    nil,
	})
}

func FindAllRoles(c *gin.Context) {
	var roles []models.Role

	if err := database.DB.Find(&roles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch roles",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Roles",
		Data:    roles,
	})
}
