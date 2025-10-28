package admin

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

	query := database.DB.Preload("user").Model(&models.Page{})
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
		})
	}

	helpers.PaginateResponse(c, pageResponses, total, page, limit, baseURL, search, "List Data Pages")
}

func CreatePage(c *gin.Context) {
	var req structs.PageCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create pages",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	username, ok := c.Get("username")
	if !ok {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "User not authenticated",
		})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", username).Find(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  map[string]string{"user": "User data not found in database"},
		})
		return
	}

	page := models.Page{
		Title:   req.Title,
		Slug:    helpers.Slugify(req.Title),
		Content: req.Content,
		UserID:  user.ID,
	}

	if err := database.DB.Create(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create page",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Page created successfully",
		Data: structs.PageResponse{
			ID:        page.ID,
			Title:     page.Title,
			Slug:      page.Slug,
			Content:   page.Content,
			UserID:    page.UserID,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindPageByID(c *gin.Context) {
	id := c.Param("id")
	var page models.Page

	if err := database.DB.First(&page, "id = ?", id).Error; err != nil {
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
		Data: structs.PageResponse{
			ID:        page.ID,
			Title:     page.Title,
			Slug:      page.Slug,
			Content:   page.Content,
			UserID:    page.UserID,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func UpdatePage(c *gin.Context) {
	id := c.Param("id")
	var req structs.PageUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var page models.Page
	if err := database.DB.First(&page, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Page not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	page.Title = req.Title
	page.Slug = helpers.Slugify(req.Title)
	page.Content = req.Content

	if err := database.DB.Save(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update page",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success update page",
		Data: structs.PageResponse{
			ID:        page.ID,
			Title:     page.Title,
			Slug:      page.Slug,
			Content:   page.Content,
			UserID:    page.UserID,
			CreatedAt: page.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: page.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func DeletePage(c *gin.Context) {
	id := c.Param("id")
	var page models.Page

	if err := database.DB.First(&page, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Page not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&page).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete page",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete post",
		Data:    nil,
	})
}
