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

func FindProducts(c *gin.Context) {
	var products []models.Product
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Preload("User").Model(&models.Product{})

	if search != "" {
		query = query.Where("title LIKE ? OR owner LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch products",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	productResponses := []structs.ProductWithRelationResponse{}
	for _, product := range products {
		productResponses = append(productResponses, structs.ProductWithRelationResponse{
			ID:      product.ID,
			Title:   product.Title,
			Slug:    product.Slug,
			Image:   product.Image,
			Owner:   product.Owner,
			Price:   product.Price,
			Phone:   product.Phone,
			Address: product.Address,
			User: structs.UserSimpleResponse{
				ID:   product.User.ID,
				Name: product.User.Name,
			},
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	helpers.PaginateResponse(c, productResponses, total, page, limit, baseURL, search, "List Data Products")
}

func CreateProduct(c *gin.Context) {
	var req structs.ProductCreateRequest

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
		DestinationDir: "public/uploads/products",
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

	product := models.Product{
		Image:   uploadResult.FileName,
		Title:   req.Title,
		Slug:    helpers.Slugify(req.Title),
		Content: req.Content,
		Owner:   req.Owner,
		Price:   req.Price,
		Phone:   req.Phone,
		Address: req.Address,
		UserID:  user.ID,
	}

	if err := database.DB.Create(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create product",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Success create product",
		Data: structs.ProductResponse{
			ID:        product.ID,
			Title:     product.Title,
			Slug:      product.Slug,
			Content:   product.Content,
			Image:     product.Image,
			Owner:     product.Owner,
			Price:     product.Price,
			Address:   product.Address,
			Phone:     product.Phone,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindProductByID(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	err := database.DB.First(&product, "id = ?", id).Error
	if err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Product not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Product found",
		Data: structs.ProductResponse{
			ID:        product.ID,
			Title:     product.Title,
			Slug:      product.Slug,
			Content:   product.Content,
			Image:     product.Image,
			Owner:     product.Owner,
			Price:     product.Price,
			Phone:     product.Phone,
			Address:   product.Address,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.DB.First(&product, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Product not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.ProductUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	oldImagePath := ""
	if product.Image != "" {
		oldImagePath = filepath.Join("public", "uploads", "products", product.Image)
	}

	file, err := c.FormFile("image")
	if err == nil {
		uploadResult := helpers.UploadFile(c, helpers.UploadConfig{
			File:           file,
			AllowedTypes:   []string{".jpg", ".jpeg", ".png", ".gif"},
			MaxSize:        10 << 20,
			DestinationDir: "public/uploads/products",
		})

		if uploadResult.Response != nil {
			c.JSON(http.StatusBadRequest, uploadResult.Response)
			return
		}

		product.Image = uploadResult.FileName
	}

	product.Title = req.Title
	product.Content = req.Content
	product.Owner = req.Owner
	product.Price = req.Price
	product.Phone = req.Phone
	product.Address = req.Address

	if err := database.DB.Save(&product).Error; err != nil {
		if file != nil && product.Image != "" {
			os.Remove(filepath.Join("public", "uploads", "products", product.Image))
		}
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update product",
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
		Data: structs.ProductResponse{
			ID:        product.ID,
			Title:     product.Title,
			Slug:      product.Slug,
			Content:   product.Content,
			Image:     product.Image,
			Owner:     product.Owner,
			Price:     product.Price,
			Address:   product.Address,
			Phone:     product.Phone,
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	if err := database.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "Product not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	imagePath := ""
	if product.Image != "" {
		imagePath = filepath.Join("public", "uploads", "products", product.Image)
	}

	if err := database.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete product",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if imagePath != "" {
		if err := os.Remove(imagePath); err != nil && !os.IsNotExist(err) {
			c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
				Success: false,
				Message: "Product deleted, but failed to delete image",
				Errors:  map[string]string{"image": err.Error()},
			})
			return
		}
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete product",
	})
}
