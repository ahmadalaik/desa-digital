package structs

type RoleCreateRequest struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []uint `json:"permission_ids"`
}

type RoleUpdateRequest struct {
	Name          string `json:"name" binding:"required"`
	PermissionIDs []uint `json:"permission_ids"`
}

type RoleResponse struct {
	ID          uint                 `json:"id"`
	Name        string               `json:"name"`
	Permissions []PermissionResponse `json:"permissions"`
	CreatedAt   string               `json:"created_at"`
	UpdatedAt   string               `json:"updated_at"`
}
