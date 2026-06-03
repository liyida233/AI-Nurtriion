package admin

import (
	"errors"
	"net/http"

	"ai-nutrition/backend/internal/httpctx"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB) {
	controller := Controller{service: NewService(db)}
	router.GET("/users", controller.ListUsers)
	router.GET("/users/:id", controller.GetUser)
	router.PATCH("/users/:id/role", controller.UpdateRole)
}

func (c Controller) ListUsers(ctx *gin.Context) {
	users, err := c.service.ListUsers(ctx.Request.Context(), httpctx.PaginationFromQuery(ctx))
	respond(ctx, users, err)
}

func (c Controller) GetUser(ctx *gin.Context) {
	user, err := c.service.GetUser(ctx.Request.Context(), ctx.Param("id"))
	respond(ctx, user, err)
}

func (c Controller) UpdateRole(ctx *gin.Context) {
	var req RoleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		httpctx.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}
	user, err := c.service.UpdateRole(ctx.Request.Context(), ctx.Param("id"), req.Role)
	respond(ctx, user, err)
}

func respond(ctx *gin.Context, data any, err error) {
	if err != nil {
		status := http.StatusInternalServerError
		message := err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
			message = "resource not found"
		}
		if message == "role must be user or admin" {
			status = http.StatusBadRequest
		}
		httpctx.Error(ctx, status, message)
		return
	}
	httpctx.OK(ctx, data)
}
