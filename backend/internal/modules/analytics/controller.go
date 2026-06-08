package analytics

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
	router.GET("/summary", controller.Summary)
	router.GET("/snapshots", controller.Snapshots)
}

func (c Controller) Summary(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	period := ctx.DefaultQuery("period", "weekly")
	startDate := ctx.Query("startDate")
	endDate := ctx.Query("endDate")
	persist := ctx.DefaultQuery("persist", "false") == "true"

	start, end, normalizedPeriod, err := ResolveRange(period, startDate, endDate)
	if err != nil {
		httpctx.Error(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if cached, ok := c.service.CachedSummary(ctx.Request.Context(), userID, normalizedPeriod); ok && !persist {
		httpctx.OK(ctx, cached)
		return
	}

	summary, err := c.service.BuildSummary(ctx.Request.Context(), userID, normalizedPeriod, start, end, persist)
	if err == nil {
		c.service.CacheSummary(ctx.Request.Context(), userID, summary)
	}
	respond(ctx, summary, err)
}

func (c Controller) Snapshots(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	snapshots, err := c.service.ListSnapshots(ctx.Request.Context(), userID)
	respond(ctx, snapshots, err)
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
