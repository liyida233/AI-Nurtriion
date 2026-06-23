package goal

import (
	"errors"
	"net/http"
	"strings"

	"ai-nutrition/backend/internal/httpctx"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	controller := Controller{service: NewService(db, redisClient)}
	router.GET("", controller.List)
	router.POST("", controller.Create)
	router.GET("/:id", controller.Get)
	router.PUT("/:id", controller.Update)
	router.PATCH("/:id/status", controller.UpdateStatus)
	router.PATCH("/:id/milestones/:milestoneId", controller.UpdateMilestone)
	router.DELETE("/:id", controller.Delete)
}

func (c Controller) List(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	goals, err := c.service.List(ctx.Request.Context(), userID)
	respond(ctx, goals, err)
}

func (c Controller) Create(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req GoalRequest
	if !bind(ctx, &req) {
		return
	}
	goal, err := c.service.Create(ctx.Request.Context(), userID, req)
	respondCreated(ctx, goal, err)
}

func (c Controller) Get(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	goal, err := c.service.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, goal, err)
}

func (c Controller) Update(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req GoalRequest
	if !bind(ctx, &req) {
		return
	}
	goal, err := c.service.Update(ctx.Request.Context(), userID, ctx.Param("id"), req)
	respond(ctx, goal, err)
}

func (c Controller) UpdateStatus(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req StatusRequest
	if !bind(ctx, &req) {
		return
	}
	goal, err := c.service.UpdateStatus(ctx.Request.Context(), userID, ctx.Param("id"), req.Status)
	respond(ctx, goal, err)
}

func (c Controller) UpdateMilestone(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req MilestoneRequest
	if !bind(ctx, &req) {
		return
	}
	milestone, err := c.service.UpdateMilestone(ctx.Request.Context(), userID, ctx.Param("id"), ctx.Param("milestoneId"), req.Completed)
	respond(ctx, milestone, err)
}

func (c Controller) Delete(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	err := c.service.Delete(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, gin.H{"deleted": true}, err)
}

func bind(ctx *gin.Context, req any) bool {
	if err := ctx.ShouldBindJSON(req); err != nil {
		httpctx.Error(ctx, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func userID(ctx *gin.Context) (string, bool) {
	value, err := httpctx.UserID(ctx)
	if err != nil {
		httpctx.Error(ctx, http.StatusUnauthorized, err.Error())
		return "", false
	}
	return value, true
}

func respond(ctx *gin.Context, data any, err error) {
	if err != nil {
		status := http.StatusInternalServerError
		message := err.Error()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			status = http.StatusNotFound
			message = "resource not found"
		} else if isValidationError(message) {
			status = http.StatusUnprocessableEntity
		}
		httpctx.Error(ctx, status, message)
		return
	}
	httpctx.OK(ctx, data)
}

func respondCreated(ctx *gin.Context, data any, err error) {
	if err != nil {
		respond(ctx, data, err)
		return
	}
	httpctx.Created(ctx, data)
}

func isValidationError(message string) bool {
	prefixes := []string{
		"deadline must be",
		"goal deadline should",
		"weight loss target appears",
		"workout frequency target should",
		"milestone dueDate must be",
		"status must be",
	}
	for _, prefix := range prefixes {
		if strings.HasPrefix(message, prefix) {
			return true
		}
	}
	return false
}
