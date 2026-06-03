package nutrition

import (
	"net/http"
	"time"

	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	redis *redis.Client
}

type createMealRequest struct {
	FoodItemID string  `json:"foodItemId" binding:"required"`
	MealType   string  `json:"mealType" binding:"required"`
	Quantity   float64 `json:"quantity" binding:"required,min=0.01"`
	MealTime   string  `json:"mealTime" binding:"required"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{db: db, redis: redisClient}
	router.GET("/foods", handler.ListFoods)
	router.POST("/foods", handler.CreateFood)
	router.GET("/meals", handler.ListMeals)
	router.POST("/meals", handler.CreateMeal)
}

func (h Handler) ListFoods(c *gin.Context) {
	var foods []models.FoodItem
	if err := h.db.Order("name asc").Find(&foods).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not list food items")
		return
	}
	httpctx.OK(c, foods)
}

func (h Handler) CreateFood(c *gin.Context) {
	var food models.FoodItem
	if err := c.ShouldBindJSON(&food); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	food.ID = uuid.NewString()
	if err := h.db.Create(&food).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not create food item")
		return
	}
	httpctx.Created(c, food)
}

func (h Handler) ListMeals(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var meals []models.MealLog
	if err := h.db.Preload("FoodItem").Where("user_id = ?", userID).Order("meal_time desc").Limit(100).Find(&meals).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not list meals")
		return
	}
	httpctx.OK(c, meals)
}

func (h Handler) CreateMeal(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var req createMealRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	mealTime, err := time.Parse(time.RFC3339, req.MealTime)
	if err != nil {
		httpctx.Error(c, http.StatusBadRequest, "mealTime must be RFC3339 datetime")
		return
	}

	var food models.FoodItem
	if err := h.db.First(&food, "id = ?", req.FoodItemID).Error; err != nil {
		httpctx.Error(c, http.StatusNotFound, "food item not found")
		return
	}

	meal := models.MealLog{
		ID:            uuid.NewString(),
		UserID:        userID,
		FoodItemID:    food.ID,
		MealType:      req.MealType,
		Quantity:      req.Quantity,
		MealTime:      mealTime,
		TotalCalories: food.Calories * req.Quantity,
		TotalProtein:  food.Protein * req.Quantity,
		TotalCarbs:    food.Carbohydrates * req.Quantity,
		TotalFat:      food.Fat * req.Quantity,
	}

	if err := h.db.Create(&meal).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not save meal")
		return
	}

	shared.InvalidateDashboard(c.Request.Context(), h.redis, userID)
	httpctx.Created(c, meal)
}
