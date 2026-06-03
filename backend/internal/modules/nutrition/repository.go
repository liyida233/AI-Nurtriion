package nutrition

import (
	"context"

	"ai-nutrition/backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) ListFoods(ctx context.Context) ([]models.FoodItem, error) {
	var foods []models.FoodItem
	err := r.db.WithContext(ctx).Order("name asc").Find(&foods).Error
	return foods, err
}

func (r Repository) GetFood(ctx context.Context, id string) (models.FoodItem, error) {
	var food models.FoodItem
	err := r.db.WithContext(ctx).First(&food, "id = ?", id).Error
	return food, err
}

func (r Repository) SaveFood(ctx context.Context, food *models.FoodItem) error {
	return r.db.WithContext(ctx).Save(food).Error
}

func (r Repository) DeleteFood(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&models.FoodItem{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

func (r Repository) ListMeals(ctx context.Context, userID string) ([]models.MealLog, error) {
	var meals []models.MealLog
	err := r.db.WithContext(ctx).
		Preload("FoodItem").
		Where("user_id = ?", userID).
		Order("meal_time desc").
		Limit(100).
		Find(&meals).Error
	return meals, err
}

func (r Repository) GetMeal(ctx context.Context, userID, id string) (models.MealLog, error) {
	var meal models.MealLog
	err := r.db.WithContext(ctx).
		Preload("FoodItem").
		Where("id = ? AND user_id = ?", id, userID).
		First(&meal).Error
	return meal, err
}

func (r Repository) SaveMeal(ctx context.Context, meal *models.MealLog) error {
	return r.db.WithContext(ctx).Save(meal).Error
}

func (r Repository) DeleteMeal(ctx context.Context, userID, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&models.MealLog{}, "id = ? AND user_id = ?", id, userID)
	return result.RowsAffected, result.Error
}
