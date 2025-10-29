package public

import (
	"net/http"

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
			Address: product.Address,
			Phone:   product.Phone,
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

func FindProductBySlug(c *gin.Context) {
	slug := c.Param("slug")
	var product models.Product

	err := database.DB.Preload("User").First(&product, "slug = ?", slug).Error
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
		Data: structs.ProductWithRelationResponse{
			ID:      product.ID,
			Title:   product.Title,
			Slug:    product.Slug,
			Content: product.Content,
			Image:   product.Image,
			Owner:   product.Owner,
			Price:   product.Price,
			Address: product.Address,
			Phone:   product.Phone,
			User: structs.UserSimpleResponse{
				ID:   product.User.ID,
				Name: product.User.Name,
			},
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		},
	})
}

func FindProductsHome(c *gin.Context) {
	var products []models.Product

	err := database.DB.Preload("User").Order("id DESC").Limit(6).Find(&products).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch posts",
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
			Address: product.Address,
			Phone:   product.Phone,
			User: structs.UserSimpleResponse{
				ID:   product.User.ID,
				Name: product.User.Name,
			},
			CreatedAt: product.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: product.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "List Data Products Home",
		Data:    productResponses,
	})
}
