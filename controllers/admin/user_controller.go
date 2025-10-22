package admin

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
)

func FindUsers(c *gin.Context) {
	var users []models.User
	var userResponses []structs.UserResponse
	var total int64

	search, page, limit, offset := helpers.GetPaginationParams(c)
	baseURL := helpers.BuildBaseURL(c)

	query := database.DB.Preload("Roles").Model(&models.User{})
	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}
	query.Count(&total)

	err := query.Order("id DESC").Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to fetch users",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	for _, user := range users {
		roleResponses := []structs.RoleResponse{}

		for _, role := range user.Roles {
			roleResponses = append(roleResponses, structs.RoleResponse{
				ID:        role.ID,
				Name:      role.Name,
				CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: role.UpdatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		userResponses = append(userResponses, structs.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Username:  user.Username,
			Email:     user.Email,
			Roles:     roleResponses,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	helpers.PaginateResponse(c, userResponses, total, page, limit, baseURL, search, "List Data Users")
}

func CreateUser(c *gin.Context) {
	var req structs.UserCreateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Error",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var roles []models.Role
	if len(req.RoleIDs) > 0 {
		database.DB.Where("id IN ?", req.RoleIDs).Find(&roles)
	}

	user := models.User{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Roles:    roles,
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to create user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusCreated, structs.SuccessResponse{
		Success: true,
		Message: "Success create user",
		Data:    user,
	})
}

func FindUserByID(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.Preload("Roles").Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	rolesResponses := []structs.RoleResponse{}
	for _, role := range user.Roles {
		rolesResponses = append(rolesResponses, structs.RoleResponse{
			ID:        role.ID,
			Name:      role.Name,
			CreatedAt: role.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: role.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	userResponse := structs.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Username:  user.Username,
		Email:     user.Email,
		Roles:     rolesResponses,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "User found",
		Data:    userResponse,
	})
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.Where("id = ?", id).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	var req structs.UserUpdateRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	hashPass, err := helpers.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Something went wrong",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	user.Name = req.Name
	user.Username = req.Username
	user.Email = req.Email
	user.Password = hashPass

	var roles []models.Role
	if len(req.RoleIDs) > 0 {
		database.DB.Where("id IN ?", req.RoleIDs).Find(&roles)
	}
	database.DB.Model(&user).Association("Roles").Replace(&roles)

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to update user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success update user",
		Data:    user,
	})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User

	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, structs.ErrorResponse{
			Success: false,
			Message: "User not found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Table("user_roles").Where("user_id = ?", id).Delete(nil).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to detach user from roles",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	if err := database.DB.Delete(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Failed to delete user",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Success delete user",
		Data:    nil,
	})
}
