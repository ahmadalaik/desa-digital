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
				ID:   post.CategoryID,
				Name: post.Category.Name,
			},
			User: structs.UserSimpleResponse{
				ID:   post.UserID,
				Name: post.User.Name,
			},
			CreatedAt: post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: post.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	helpers.PaginateResponse(c, postResponses, total, page, limit, baseURL, search, "List Data Posts")
}

func CreatePost(c *gin.Context) {
	var req structs.PostCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create post",
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
		DestinationDir: "public/uploads/posts",
	})
	if uploadResult.Response != nil {
		c.JSON(http.StatusBadRequest, uploadResult.Response)
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
	if err := database.DB.Where("username = ?", username).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  map[string]string{"user": "User data not found in database"},
		})
		return
	}

	post := models.Post{
		Title:      req.Title,
		Slug:       helpers.Slugify(req.Title),
		Content:    req.Content,
		Image:      uploadResult.FileName,
		CategoryID: req.CategoryID,
		UserID:     user.ID,
	}

	if err := database.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create post",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Post created successfully",
		Data: structs.PostResponse{
			ID:         post.ID,
			Image:      post.Image,
			Title:      post.Title,
			Slug:       post.Slug,
			Content:    post.Content,
			CategoryID: post.CategoryID,
			UserID:     post.UserID,
			CreatedAt:  post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindPostByID(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	err := database.DB.Preload("Category").Preload("User").First(&post, "id = ?", id).Error
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
		Data: structs.PostResponse{
			ID:         post.ID,
			Image:      post.Image,
			Title:      post.Title,
			Slug:       post.Slug,
			Content:    post.Content,
			CategoryID: post.CategoryID,
			UserID:     post.UserID,
			CreatedAt:  post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func UpdatePost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := database.DB.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Post not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.PostUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	oldImagePath := ""
	if post.Image != "" {
		oldImagePath = filepath.Join("public", "uploads", "posts", post.Image)
	}

	file, err := c.FormFile("image")
	if err == nil {
		uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
			File:           file,
			AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".gif"},
			MaxSize:        10 << 20,
			DestinationDir: "public/uploads/posts",
		})

		if uploadResult.Response != nil {
			c.JSON(http.StatusBadRequest, uploadResult.Response)
			return
		}

		post.Image = uploadResult.FileName
	}

	post.Title = req.Title
	post.Slug = helpers.Slugify(req.Title)
	post.Content = req.Content
	post.CategoryID = req.CategoryID

	if err := database.DB.Save(&post).Error; err != nil {
		if file != nil && post.Image != "" {
			newImagePath := filepath.Join("public", "uploads", "posts", post.Image)
			os.Remove(newImagePath)
		}

		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update post",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if file != nil && oldImagePath != "" {
		os.Remove(oldImagePath)
	}

	database.DB.Preload("Category").Preload("User").First(&post, post.ID)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success update post",
		Data: structs.PostResponse{
			ID:         post.ID,
			Image:      post.Image,
			Title:      post.Title,
			Slug:       post.Slug,
			Content:    post.Content,
			CategoryID: post.CategoryID,
			UserID:     post.UserID,
			CreatedAt:  post.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:  post.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func DeletPost(c *gin.Context) {
	id := c.Param("id")
	var post models.Post

	if err := database.DB.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Post not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	imagePath := ""
	if post.Image != "" {
		imagePath = filepath.Join("public", "uploads", "posts", post.Image)
	}

	if err := database.DB.Delete(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete post",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Post deleted but failed to remove image",
				Errors:  map[string]string{"image": "Failed to remove image file: " + err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete post",
		Data: nil,
	})
}
