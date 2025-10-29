package public

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindAparaturs(c *gin.Context) {
	var aparaturs []models.Aparatur
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Aparatur{})
	if search != "" {
		query = query.Where("name LIKE ? OR position LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&aparaturs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch aparaturs",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	helpers.PaginateResponse(c, aparaturs, total, page, limit, baseURL, search, "List Data Aparaturs")
}

func FindAparaturByID(c *gin.Context) {
	id := c.Param("id")
	var aparatur models.Aparatur

	if err := database.DB.First(&aparatur, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Aparatur not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Aparatur found",
		Data: structs.AparaturResponse{
			ID:          aparatur.ID,
			Image:       aparatur.Image,
			Name:        aparatur.Name,
			Position:    aparatur.Position,
			Description: aparatur.Description,
			CreatedAt:   aparatur.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:   aparatur.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindAparatursHome(c *gin.Context) {
	var aparaturs []models.Aparatur

	err := database.DB.Order("id DESC").Limit(6).Find(&aparaturs).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch aparaturs",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Aparatur found",
		Data:    aparaturs,
	})
}
