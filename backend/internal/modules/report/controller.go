package report

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

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, _ *redis.Client) {
	controller := Controller{service: NewService(db)}
	router.GET("", controller.List)
	router.POST("/generate", controller.Generate)
	router.GET("/:id/download", controller.Download)
	router.GET("/:id", controller.Get)
	router.DELETE("/:id", controller.Delete)
}

func (c Controller) List(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	reports, err := c.service.List(ctx.Request.Context(), userID)
	respond(ctx, reports, err)
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
	report, err := c.service.Generate(ctx.Request.Context(), userID, req)
	respondCreated(ctx, report, err)
}

func (c Controller) Get(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	report, err := c.service.Get(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, report, err)
}

func (c Controller) Delete(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	err := c.service.Delete(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, gin.H{"deleted": true}, err)
}

func (c Controller) Download(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	path, err := c.service.ReportFilePath(ctx.Request.Context(), userID, ctx.Param("id"))
	if err != nil {
		respond(ctx, nil, err)
		return
	}
	ctx.FileAttachment(path, "progress-report-"+ctx.Param("id")+".html")
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
		if message == "periodType must be weekly or monthly" {
			status = http.StatusBadRequest
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
