package auth

import (
	"net/http"

	"github.com/ahmadalaik/desa-digital/database"
	"github.com/ahmadalaik/desa-digital/helpers"
	"github.com/ahmadalaik/desa-digital/models"
	"github.com/ahmadalaik/desa-digital/structs"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context) {
	var req structs.UserLoginRequest
	var user models.User

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, structs.ErrorResponse{
			Success: false,
			Message: "Validation Errors",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	err = database.DB.Preload("Roles").Preload("Roles.Permissions").Where("username = ?", req.Username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "User Not Found",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, structs.ErrorResponse{
			Success: false,
			Message: "Username or Password is Wrong",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	token, err := helpers.GenerateToken(user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, structs.ErrorResponse{
			Success: false,
			Message: "Someting went wrong",
			Errors:  helpers.TranslateErrorMessage(err),
		})
		return
	}

	permissionMap := helpers.GetPermission(user.Roles)

	c.JSON(http.StatusOK, structs.SuccessResponse{
		Success: true,
		Message: "Login Success",
		Data: structs.UserResponse{
			ID:          user.ID,
			Name:        user.Name,
			Username:    user.Username,
			Email:       user.Email,
			Permissions: permissionMap,
			Token:       &token,
		},
	})
}
