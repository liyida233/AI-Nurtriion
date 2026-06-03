package recommendation

import (
	"errors"
	"net/http"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/httpctx"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Controller struct {
	service Service
}

func RegisterRoutes(router *gin.RouterGroup, cfg config.Config, db *gorm.DB, redisClient *redis.Client) {
	controller := Controller{service: NewService(cfg, db, redisClient)}
	router.GET("", controller.List)
	router.POST("/generate", controller.Generate)
	router.GET("/:id", controller.Get)
	router.DELETE("/:id", controller.Delete)
	router.GET("/:id/feedback", controller.ListFeedback)
	router.POST("/:id/feedback", controller.AddFeedback)
}

func (c Controller) List(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	recommendations, err := c.service.List(ctx.Request.Context(), userID)
	respond(ctx, recommendations, err)
}

func (c Controller) Generate(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req GenerateRequest
	if !bind(ctx, &req) {
		return
	}
	recommendation, err := c.service.Generate(ctx.Request.Context(), userID, req)
	respondCreated(ctx, recommendation, err)
}

func (c Controller) Get(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	recommendation, err := c.service.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, recommendation, err)
}

func (c Controller) Delete(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	err := c.service.Delete(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, gin.H{"deleted": true}, err)
}

func (c Controller) AddFeedback(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req FeedbackRequest
	if !bind(ctx, &req) {
		return
	}
	feedback, err := c.service.AddFeedback(ctx.Request.Context(), userID, ctx.Param("id"), req)
	respondCreated(ctx, feedback, err)
}

func (c Controller) ListFeedback(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	feedback, err := c.service.ListFeedback(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, feedback, err)
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
		}
		if message == "recommendation did not pass safety validation" {
			status = http.StatusUnprocessableEntity
		}
		if message == "AI recommendation rate limit exceeded" {
			status = http.StatusTooManyRequests
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
