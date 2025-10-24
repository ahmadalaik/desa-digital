package admin

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindCategories(c *gin.Context) {
	var categories []models.Category
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Category{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&categories).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch categories",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	helpers.PaginateResponse(c, categories, total, page, limit, baseURL, search, "List Data Categories")
}

func CreateCategory(c *gin.Context) {
	var req structs.CategoryCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	category := models.Category{
		Name: req.Name,
	}

	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create category",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Success create category",
		Data:    category,
	})
}

func FindCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.DB.First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Category found",
		Data:    category,
	})
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.DB.First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.CategoryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	category.Name = req.Name
	category.Slug = helpers.Slugify(req.Name)

	if err := database.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update category",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success update category",
		Data:    category,
	})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	var category models.Category

	if err := database.DB.First(&category, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Category not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete category",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}
}

func FindAllCategories(c *gin.Context) {
	var categories []models.Category

	if err := database.DB.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch categories",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Categories",
		Data:    categories,
	})
}
