package structs

type (
	PageCreateRequest struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}

	PageUpdateRequest struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
)

type (
	PageResponse struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Slug      string `json:"slug"`
		Content   string `json:"content"`
		UserID    uint   `json:"user_id"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	PageWithRelationResponse struct {
		ID        uint               `json:"id"`
		Title     string             `json:"title"`
		Slug      string             `json:"slug"`
		Content   string             `json:"content,omitempty"`
		User      UserSimpleResponse `json:"user,omitempty"`
		CreatedAt string             `json:"created_at"`
		UpdatedAt string             `json:"updated_at"`
	}
)
