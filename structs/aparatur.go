package structs

type (
	AparaturCreateRequest struct {
		Name        string `json:"name" binding:"required"`
		Position    string `json:"position" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	AparaturUpdateRequest struct {
		Name        string `json:"name" binding:"required"`
		Position    string `json:"position" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
)

type (
	AparaturResponse struct {
		ID          uint   `json:"id"`
		Image       string `json:"image"`
		Name        string `json:"name"`
		Position    string `json:"position"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
		UpdatedAt   string `json:"updated_at"`
	}
)
