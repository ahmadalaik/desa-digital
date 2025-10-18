package structs

type (
	CategoryCreateRequest struct {
		Name string `json:"name" binding:"required"`
	}

	CategoryUpdateRequest struct {
		Name string `json:"name" binding:"required"`
	}
)

type (
	CategoryResponse struct {
		ID        uint   `json:"id"`
		Name      string `json:"name"`
		Slug      string `json:"slug"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	CategorySimpleResponse struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}
)
