package public

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindSliders(c *gin.Context) {
	var sliders []models.Slider

	if err := database.DB.Find(&sliders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch sliders",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Lists Data Sliders",
		Data:    sliders,
	})
}
