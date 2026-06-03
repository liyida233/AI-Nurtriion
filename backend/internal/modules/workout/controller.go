package workout

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
	router.GET("/exercises", controller.ListExercises)
	router.GET("/exercises/:id", controller.GetExercise)
	router.GET("", controller.ListSessions)
	router.POST("", controller.CreateSession)
	router.GET("/:id", controller.GetSession)
	router.PUT("/:id", controller.UpdateSession)
	router.DELETE("/:id", controller.DeleteSession)
}

func RegisterAdminRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	controller := Controller{service: NewService(db, redisClient)}
	router.POST("/exercises", controller.CreateExercise)
	router.PUT("/exercises/:id", controller.UpdateExercise)
	router.DELETE("/exercises/:id", controller.DeleteExercise)
}

func (c Controller) ListExercises(ctx *gin.Context) {
	exercises, err := c.service.ListExercises(ctx.Request.Context())
	respond(ctx, exercises, err)
}

func (c Controller) CreateExercise(ctx *gin.Context) {
	var req ExerciseRequest
	if !bind(ctx, &req) {
		return
	}
	exercise, err := c.service.CreateExercise(ctx.Request.Context(), req)
	respondCreated(ctx, exercise, err)
}

func (c Controller) GetExercise(ctx *gin.Context) {
	exercise, err := c.service.GetExercise(ctx.Request.Context(), ctx.Param("id"))
	respond(ctx, exercise, err)
}

func (c Controller) UpdateExercise(ctx *gin.Context) {
	var req ExerciseRequest
	if !bind(ctx, &req) {
		return
	}
	exercise, err := c.service.UpdateExercise(ctx.Request.Context(), ctx.Param("id"), req)
	respond(ctx, exercise, err)
}

func (c Controller) DeleteExercise(ctx *gin.Context) {
	err := c.service.DeleteExercise(ctx.Request.Context(), ctx.Param("id"))
	respond(ctx, gin.H{"deleted": true}, err)
}

func (c Controller) ListSessions(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	sessions, err := c.service.ListSessions(ctx.Request.Context(), userID)
	respond(ctx, sessions, err)
}

func (c Controller) CreateSession(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req SessionRequest
	if !bind(ctx, &req) {
		return
	}
	session, indicators, err := c.service.CreateSession(ctx.Request.Context(), userID, req)
	respondCreated(ctx, gin.H{"workout": session, "indicators": indicators}, err)
}

func (c Controller) GetSession(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	session, err := c.service.GetSession(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, session, err)
}

func (c Controller) UpdateSession(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req SessionRequest
	if !bind(ctx, &req) {
		return
	}
	session, indicators, err := c.service.UpdateSession(ctx.Request.Context(), userID, ctx.Param("id"), req)
	respond(ctx, gin.H{"workout": session, "indicators": indicators}, err)
}

func (c Controller) DeleteSession(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	err := c.service.DeleteSession(ctx.Request.Context(), userID, ctx.Param("id"))
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
