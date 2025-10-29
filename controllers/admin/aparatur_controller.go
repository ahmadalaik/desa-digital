package admin

import (
	"net/http"
	"os"
	"path/filepath"

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

func CreateAparatur(c *gin.Context) {
	var req structs.AparaturCreateRequest

	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  map[string]string{"Image": "Image is required"},
		})
		return
	}

	uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
		File:           file,
		AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".gif"},
		MaxSize:        10 << 20,
		DestinationDir: "public/uploads/aparaturs",
	})

	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
		return
	}

	aparatur := models.Aparatur{
		Name:        req.Name,
		Position:    req.Position,
		Description: req.Description,
		Image:       uploadResult.FileName,
	}

	if err := database.DB.Create(&aparatur).Error; err != nil {
		if uploadResult.FileName != "" {
			os.Remove(filepath.Join("public", "uploads", "aparaturs", uploadResult.FileName))
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create aparatur",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Success create aparatur",
		Data:    aparatur,
	})
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

func UpdateAparatur(c *gin.Context) {
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

	var req structs.AparaturUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	oldImagePath := ""
	if aparatur.Image != "" {
		oldImagePath = filepath.Join("public", "uploads", "aparaturs", aparatur.Image)
	}

	file, err := c.FormFile("image")
	if err == nil {
		uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
			File:           file,
			AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".gif"},
			MaxSize:        10 << 20,
			DestinationDir: "public/uploads/aparaturs",
		})

		if uploadResult.Response != nil {
			c.JSON(http.StatusBadRequest, uploadResult.Response)
			return
		}

		aparatur.Image = uploadResult.FileName
	}

	aparatur.Name = req.Name
	aparatur.Position = req.Position
	aparatur.Description = req.Description

	if err := database.DB.Save(&aparatur).Error; err != nil {
		if file != nil && aparatur.Image != "" {
			os.Remove(filepath.Join("public", "uploads", "aparaturs", aparatur.Image))
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update aparatur",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if file != nil && oldImagePath != "" {
		os.Remove(oldImagePath)
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success update product",
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

func DeleteAparatur(c *gin.Context) {
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

	imagePath := ""
	if aparatur.Image != "" {
		imagePath = filepath.Join("public", "uploads", "aparaturs", aparatur.Image)
	}

	if err := database.DB.Delete(&aparatur).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete aparatur",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Aparatur deleted but failed to remove image",
				Errors:  map[string]string{"image": err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete aparatur",
	})
}
