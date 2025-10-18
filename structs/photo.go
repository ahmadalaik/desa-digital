package structs

type (
	PhotoCreateRequest struct {
		Caption     string `json:"caption" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
)

type (
	PhotoResponse struct {
		ID          uint   `json:"id"`
		Image       string `json:"image"`
		Caption     string `json:"caption"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
)
