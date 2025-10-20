package admin

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindPermissons(c *gin.Context) {
	var permissions []models.Permission
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Permission{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&permissions).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch permissions",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	helpers.PaginateResponse(c, permissions, total, page, limit, baseURL, search, "List Data Permissions")
}

func CreatePermission(c *gin.Context) {
	var req structs.PermissionCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	permission := models.Permission{
		Name: req.Name,
	}

	if err := database.DB.Create(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create permission",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Success create permission",
		Data:    permission,
	})
}

func FindPermissonByID(c *gin.Context) {
	id := c.Param("id")

	var permission models.Permission
	if err := database.DB.Where("id = ?", id).First(&permission).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Permission found",
		Data:    permission,
	})
}

func UpdatePermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.Permission

	if err := database.DB.First(&permission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.PermissionUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	permission.Name = req.Name

	if err := database.DB.Save(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update permission",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success update permission",
		Data:    permission,
	})
}

func DeletePermission(c *gin.Context) {
	id := c.Param("id")
	var permission models.Permission

	if err := database.DB.First(&permission, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Permission not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Table("role_permissions").Where("permission_id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to detach permission from roles",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&permission).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete permission",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete permission",
		Data:    nil,
	})
}

func FindAllPermissions(c *gin.Context) {
	var permissions []models.Permission

	if err := database.DB.Find(&permissions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch permission",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Permissions",
		Data:    permissions,
	})
}
