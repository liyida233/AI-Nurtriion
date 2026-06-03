package body

import (
	"errors"
	"net/http"

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
	router.DELETE("/:id", controller.Delete)
}

func (c Controller) List(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	records, err := c.service.List(ctx.Request.Context(), userID)
	respond(ctx, records, err)
}

func (c Controller) Create(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req RecordRequest
	if !bind(ctx, &req) {
		return
	}
	record, err := c.service.Create(ctx.Request.Context(), userID, req)
	respondCreated(ctx, record, err)
}

func (c Controller) Get(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	record, err := c.service.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, record, err)
}

func (c Controller) Update(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req RecordRequest
	if !bind(ctx, &req) {
		return
	}
	record, err := c.service.Update(ctx.Request.Context(), userID, ctx.Param("id"), req)
	respond(ctx, record, err)
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
