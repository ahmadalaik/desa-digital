package structs

type (
	ProductCreateRequest struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Owner   string `json:"owner" binding:"required"`
		Price   int    `json:"price" binding:"required"`
		Phone   string `json:"phone" binding:"required"`
		Address string `json:"address" binding:"required"`
	}

	ProductUpdateRequest struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Owner   string `json:"owner" binding:"required"`
		Price   int    `json:"price" binding:"required"`
		Phone   string `json:"phone" binding:"required"`
		Address string `json:"address" binding:"required"`
	}
)

type (
	ProductResponse struct {
		ID        uint   `json:"id"`
		Title     string `json:"title"`
		Slug      string `json:"slug"`
		Content   string `json:"content"`
		Image     string `json:"image"`
		Owner     string `json:"owner"`
		Price     int    `json:"price"`
		Phone     string `json:"phone"`
		Address   string `json:"address"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	ProductWithRelationResponse struct {
		ID        uint               `json:"id"`
		Title     string             `json:"title"`
		Slug      string             `json:"slug"`
		Content   string             `json:"content,omitempty"`
		Image     string             `json:"image,omitempty"`
		Owner     string             `json:"owner"`
		Price     int                `json:"price,omitempty"`
		Phone     string             `json:"phone,omitempty"`
		Address   string             `json:"address,omitempty"`
		User      UserSimpleResponse `json:"user,omitempty"`
		CreatedAt string             `json:"created_at"`
		UpdatedAt string             `json:"updated_at"`
	}
)
