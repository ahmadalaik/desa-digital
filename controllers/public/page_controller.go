package public

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindPages(c *gin.Context) {
	var pages []models.Page
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Preload("User").Model(&models.Page{})
	if search != "" {
		query = query.Where("title LIKE ? OR content LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&pages).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch pages",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	pageResponses := []structs.PageWithRelationResponse{}
	for _, page := range pages {
		pageResponses = append(pageResponses, structs.PageWithRelationResponse{
			ID:    page.ID,
			Title: page.Title,
			Slug:  page.Slug,
			User: structs.UserSimpleResponse{
				ID:   page.User.ID,
				Name: page.User.Name,
			},
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	helpers.PaginateResponse(c, pageResponses, total, page, limit, baseURL, search, "List Data Pages")
}

func FindPageBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var page models.Page

	err := database.DB.Preload("User").First(&page, "slug = ?", slug).Error
	if err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Page not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Page found",
		Data: structs.PageWithRelationResponse{
			ID:      page.ID,
			Title:   page.Title,
			Slug:    page.Slug,
			Content: page.Content,
			User: structs.UserSimpleResponse{
				ID:   page.User.ID,
				Name: page.User.Name,
			},
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}
