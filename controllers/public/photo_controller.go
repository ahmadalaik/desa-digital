package public

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindPhotos(c *gin.Context) {
	var photos []models.Photo
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Photo{})
	if search != "" {
		query = query.Where("caption LIKE ?", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&photos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch photos",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	helpers.PaginateResponse(c, photos, total, page, limit, baseURL, search, "List Data Photos")
}

func FindPhotosHome(c *gin.Context) {
	var photos []models.Photo

	err := database.DB.Order("id desc").Limit(6).Find(&photos).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch photos",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Data Photos Home",
		Data:    photos,
	})
}
