package nutrition

import (
	"context"
	"errors"
	"time"

	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/shared"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	repo  Repository
	redis *redis.Client
}

func NewService(db *gorm.DB, redisClient *redis.Client) Service {
	return Service{repo: NewRepository(db), redis: redisClient}
}

func (s Service) ListFoods(ctx context.Context) ([]models.FoodItem, error) {
	return s.repo.ListFoods(ctx)
}

func (s Service) GetFood(ctx context.Context, id string) (models.FoodItem, error) {
	return s.repo.GetFood(ctx, id)
}

func (s Service) CreateFood(ctx context.Context, req FoodRequest) (models.FoodItem, error) {
	food := mapFood(models.FoodItem{ID: uuid.NewString()}, req)
	return food, s.repo.SaveFood(ctx, &food)
}

func (s Service) UpdateFood(ctx context.Context, id string, req FoodRequest) (models.FoodItem, error) {
	food, err := s.repo.GetFood(ctx, id)
	if err != nil {
		return food, err
	}
	food = mapFood(food, req)
	return food, s.repo.SaveFood(ctx, &food)
}

func (s Service) DeleteFood(ctx context.Context, id string) error {
	rows, err := s.repo.DeleteFood(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s Service) ListMeals(ctx context.Context, userID string) ([]models.MealLog, error) {
	return s.repo.ListMeals(ctx, userID)
}

func (s Service) GetMeal(ctx context.Context, userID, id string) (models.MealLog, error) {
	return s.repo.GetMeal(ctx, userID, id)
}

func (s Service) CreateMeal(ctx context.Context, userID string, req MealRequest) (models.MealLog, MealIndicators, error) {
	mealTime, err := parseMealTime(req.MealTime)
	if err != nil {
		return models.MealLog{}, MealIndicators{}, err
	}
	food, err := s.repo.GetFood(ctx, req.FoodItemID)
	if err != nil {
		return models.MealLog{}, MealIndicators{}, err
	}
	totals := CalculateMealTotals(food, req.Quantity)
	meal := models.MealLog{
		ID:            uuid.NewString(),
		UserID:        userID,
		FoodItemID:    food.ID,
		MealType:      req.MealType,
		Quantity:      req.Quantity,
		MealTime:      mealTime,
		TotalCalories: totals.Calories,
		TotalProtein:  totals.Protein,
		TotalCarbs:    totals.Carbs,
		TotalFat:      totals.Fat,
	}
	if err := s.repo.SaveMeal(ctx, &meal); err != nil {
		return meal, totals, err
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return meal, totals, nil
}

func (s Service) UpdateMeal(ctx context.Context, userID, id string, req MealRequest) (models.MealLog, MealIndicators, error) {
	mealTime, err := parseMealTime(req.MealTime)
	if err != nil {
		return models.MealLog{}, MealIndicators{}, err
	}
	meal, err := s.repo.GetMeal(ctx, userID, id)
	if err != nil {
		return meal, MealIndicators{}, err
	}
	food, err := s.repo.GetFood(ctx, req.FoodItemID)
	if err != nil {
		return meal, MealIndicators{}, err
	}
	totals := CalculateMealTotals(food, req.Quantity)
	meal.FoodItemID = food.ID
	meal.MealType = req.MealType
	meal.Quantity = req.Quantity
	meal.MealTime = mealTime
	meal.TotalCalories = totals.Calories
	meal.TotalProtein = totals.Protein
	meal.TotalCarbs = totals.Carbs
	meal.TotalFat = totals.Fat
	if err := s.repo.SaveMeal(ctx, &meal); err != nil {
		return meal, totals, err
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return meal, totals, nil
}

func (s Service) DeleteMeal(ctx context.Context, userID, id string) error {
	rows, err := s.repo.DeleteMeal(ctx, userID, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return nil
}

func mapFood(food models.FoodItem, req FoodRequest) models.FoodItem {
	food.Name = req.Name
	food.ServingSize = req.ServingSize
	food.Calories = req.Calories
	food.Protein = req.Protein
	food.Carbohydrates = req.Carbohydrates
	food.Fat = req.Fat
	food.Sugar = req.Sugar
	food.Sodium = req.Sodium
	return food
}

func parseMealTime(value string) (time.Time, error) {
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return time.Time{}, errors.New("mealTime must be RFC3339 datetime")
	}
	return parsed, nil
}
