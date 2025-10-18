package structs

type (
	UserCreateRequest struct {
		Name     string `json:"name" binding:"required"`
		Username string `json:"username" binding:"required" gorm:"unique;not null"`
		Email    string `json:"email" binding:"required" gorm:"unique;not null"`
		Password string `json:"password" binding:"required"`
		RoleIDs  []uint `json:"role_ids" binding:"required"`
	}

	UserUpdateRequest struct {
		Name     string `json:"name" binding:"required"`
		Username string `json:"username" binding:"required" gorm:"unique;not null"`
		Email    string `json:"email" binding:"required" gorm:"unique;not null"`
		Password string `json:"password,omitempty"`
		RoleIDs  []uint `json:"role_ids"`
	}

	UserLoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
)

type (
	UserResponse struct {
		ID          uint            `json:"id"`
		Name        string          `json:"name"`
		Username    string          `json:"username"`
		Email       string          `json:"email"`
		Permissions map[string]bool `json:"permissions,omitempty"`
		Roles       []RoleResponse  `json:"roles,omitempty"`
		Token       *string         `json:"token,omitempty"`
		CreatedAt   string          `json:"created_at,omitempty"`
		UpdatedAt   string          `json:"updated_at,omitempty"`
	}

	UserSimpleResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
)
