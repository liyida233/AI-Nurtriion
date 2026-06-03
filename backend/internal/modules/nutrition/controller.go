package nutrition

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
	router.GET("/foods", controller.ListFoods)
	router.GET("/foods/:id", controller.GetFood)
	router.GET("/meals", controller.ListMeals)
	router.POST("/meals", controller.CreateMeal)
	router.GET("/meals/:id", controller.GetMeal)
	router.PUT("/meals/:id", controller.UpdateMeal)
	router.DELETE("/meals/:id", controller.DeleteMeal)
}

func RegisterAdminRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	controller := Controller{service: NewService(db, redisClient)}
	router.POST("/foods", controller.CreateFood)
	router.PUT("/foods/:id", controller.UpdateFood)
	router.DELETE("/foods/:id", controller.DeleteFood)
}

func (c Controller) ListFoods(ctx *gin.Context) {
	foods, err := c.service.ListFoods(ctx.Request.Context())
	respond(ctx, foods, err)
}

func (c Controller) CreateFood(ctx *gin.Context) {
	var req FoodRequest
	if !bind(ctx, &req) {
		return
	}
	food, err := c.service.CreateFood(ctx.Request.Context(), req)
	respondCreated(ctx, food, err)
}

func (c Controller) GetFood(ctx *gin.Context) {
	food, err := c.service.GetFood(ctx.Request.Context(), ctx.Param("id"))
	respond(ctx, food, err)
}

func (c Controller) UpdateFood(ctx *gin.Context) {
	var req FoodRequest
	if !bind(ctx, &req) {
		return
	}
	food, err := c.service.UpdateFood(ctx.Request.Context(), ctx.Param("id"), req)
	respond(ctx, food, err)
}

func (c Controller) DeleteFood(ctx *gin.Context) {
	err := c.service.DeleteFood(ctx.Request.Context(), ctx.Param("id"))
	respond(ctx, gin.H{"deleted": true}, err)
}

func (c Controller) ListMeals(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	meals, err := c.service.ListMeals(ctx.Request.Context(), userID)
	respond(ctx, meals, err)
}

func (c Controller) CreateMeal(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req MealRequest
	if !bind(ctx, &req) {
		return
	}
	meal, totals, err := c.service.CreateMeal(ctx.Request.Context(), userID, req)
	respondCreated(ctx, gin.H{"meal": meal, "indicators": totals}, err)
}

func (c Controller) GetMeal(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	meal, err := c.service.GetMeal(ctx.Request.Context(), userID, ctx.Param("id"))
	respond(ctx, meal, err)
}

func (c Controller) UpdateMeal(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	var req MealRequest
	if !bind(ctx, &req) {
		return
	}
	meal, totals, err := c.service.UpdateMeal(ctx.Request.Context(), userID, ctx.Param("id"), req)
	respond(ctx, gin.H{"meal": meal, "indicators": totals}, err)
}

func (c Controller) DeleteMeal(ctx *gin.Context) {
	userID, ok := userID(ctx)
	if !ok {
		return
	}
	err := c.service.DeleteMeal(ctx.Request.Context(), userID, ctx.Param("id"))
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
