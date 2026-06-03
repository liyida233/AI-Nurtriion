package admin

type RoleRequest struct {
	Role string `json:"role" binding:"required"`
}
