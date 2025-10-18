package structs

type (
	SliderCreateRequest struct {
		Description string `json:"description" binding:"required"`
	}
)

type (
	SliderResponse struct {
		ID          uint   `json:"id"`
		Image       string `json:"image"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
)
