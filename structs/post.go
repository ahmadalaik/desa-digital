package structs

type (
	PostCreateRequest struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}

	PostUpdateRequest struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}
)

type (
	PostResponse struct {
		ID         uint   `json:"id"`
		Image      string `json:"image"`
		Title      string `json:"title"`
		Slug       string `json:"slug"`
		Content    string `json:"content"`
		CategoryID uint   `json:"category_id"`
		UserID     uint   `json:"user_id"`
		CreatedAt  string `json:"created_at"`
		UpdatedAt  string `json:"updated_at"`
	}

	PostWithRelationResponse struct {
		ID        uint                   `json:"id"`
		Image     string                 `json:"image"`
		Title     string                 `json:"title"`
		Slug      string                 `json:"slug"`
		Content   string                 `json:"content,omitempty"`
		Category  CategorySimpleResponse `json:"category,omitempty"`
		User      UserSimpleResponse     `json:"user,omitempty"`
		CreatedAt string                 `json:"created_at"`
		UpdatedAt string                 `json:"updated_at"`
	}
)
