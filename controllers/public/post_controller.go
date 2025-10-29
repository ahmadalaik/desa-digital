package public

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindPosts(c *gin.Context) {
	var posts []models.Post
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Preload("Category").Preload("User").Model(&models.Post{})
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch posts",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	postResponses := []structs.PostWithRelationResponse{}
	for _, post := range posts {
		postResponses = append(postResponses, structs.PostWithRelationResponse{
			ID:      post.ID,
			Image:   post.Image,
			Title:   post.Title,
			Slug:    post.Slug,
			Content: post.Content,
			Category: structs.CategorySimpleResponse{
				ID:   post.Category.ID,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.User.ID,
				Name: post.User.Name,
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	helpers.PaginateResponse(c, postResponses, total, page, limit, baseURL, search, "List Data Posts")
}

func FindPostBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var post models.Post

	err := database.DB.Preload("Category").Preload("User").First(&post, "slug = ?", slug).Error
	if err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Post not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Post found",
		Data: structs.PostWithRelationResponse{
			ID:      post.ID,
			Image:   post.Image,
			Title:   post.Title,
			Slug:    post.Slug,
			Content: post.Content,
			Category: structs.CategorySimpleResponse{
				ID:   post.Category.ID,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.User.ID,
				Name: post.User.Name,
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindPostsHome(c *gin.Context) {
	var posts []models.Post

	err := database.DB.Preload("Category").Preload("User").Order("id DESC").Limit(6).Find(&posts).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch posts",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	postResponses := []structs.PostWithRelationResponse{}
	for _, post := range posts {
		postResponses = append(postResponses, structs.PostWithRelationResponse{
			ID:      post.ID,
			Image:   post.Image,
			Title:   post.Title,
			Slug:    post.Slug,
			Content: post.Content,
			Category: structs.CategorySimpleResponse{
				ID:   post.Category.ID,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.User.ID,
				Name: post.User.Name,
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Data Posts Home",
		Data:    postResponses,
	})
}
