package admin

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func Dashboard(c *gin.Context) {
	var (
		categoriesCount int64
		postsCount      int64
		productsCount   int64
		aparatursCount  int64
	)

	if err := database.DB.Model(&models.Category{}).Count(&categoriesCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get categories count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Model(&models.Post{}).Count(&postsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get posts count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Model(&models.Product{}).Count(&productsCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get products count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Model(&models.Aparatur{}).Count(&aparatursCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to get aparaturs count",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Dashboard stats retrieved successfully",
		Data: structs.DashboardResponse{
			CategoriesCount: categoriesCount,
			PostsCount:      postsCount,
			ProductsCount:   productsCount,
			AparatursCount:  aparatursCount,
		},
	})
}
