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

func FindSliders(c *gin.Context) {
	var sliders []models.Slider
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Model(&models.Slider{})
	if search != "" {
		query = query.Where("description LIKE ?", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id desc").Limit(limit).Offset(offset).Find(&sliders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch sliders",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	helpers.PaginateResponse(c, sliders, total, page, limit, baseURL, search, "List Data Sliders")
}

func CreateSlider(c *gin.Context) {
	var req structs.SliderCreateRequest

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
		DestinationDir: "public/uploads/sliders",
	})

	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
		return
	}

	slider := models.Slider{
		Image:       uploadResult.FileName,
		Description: req.Description,
	}

	if err := database.DB.Create(&slider).Error; err != nil {
		if uploadResult.FileName != "" {
			os.Remove(filepath.Join("public", "uploads", "sliders", uploadResult.FileName))
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create slider",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Success create slider",
		Data:    slider,
	})
}

func DeleteSlider(c *gin.Context) {
	id := c.Param("id")
	var slider models.Slider

	if err := database.DB.First(&slider, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Slider not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	imagePath := ""
	if slider.Image != "" {
		imagePath = filepath.Join("public", "uploads", "sliders", slider.Image)
	}

	if err := database.DB.Delete(&slider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete slider",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Slider deleted but failed to remove image",
				Errors:  map[string]string{"image": err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete slider",
	})
}
