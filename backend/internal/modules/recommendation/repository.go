package recommendation

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

func (r Repository) List(ctx context.Context, userID string) ([]models.AIRecommendation, error) {
	var recommendations []models.AIRecommendation
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Limit(50).Find(&recommendations).Error
	return recommendations, err
}

func (r Repository) Get(ctx context.Context, userID, id string) (models.AIRecommendation, error) {
	var recommendation models.AIRecommendation
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&recommendation).Error
	return recommendation, err
}

func (r Repository) Save(ctx context.Context, recommendation *models.AIRecommendation) error {
	return r.db.WithContext(ctx).Save(recommendation).Error
}

func (r Repository) Delete(ctx context.Context, userID, id string) (int64, error) {
	var rows int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.RecommendationFeedback{}, "recommendation_id = ?", id).Error; err != nil {
			return err
		}
		result := tx.Delete(&models.AIRecommendation{}, "id = ? AND user_id = ?", id, userID)
		rows = result.RowsAffected
		return result.Error
	})
	return rows, err
}

func (r Repository) CreateFeedback(ctx context.Context, feedback *models.RecommendationFeedback) error {
	return r.db.WithContext(ctx).Create(feedback).Error
}

func (r Repository) ListFeedback(ctx context.Context, recommendationID string) ([]models.RecommendationFeedback, error) {
	var feedback []models.RecommendationFeedback
	err := r.db.WithContext(ctx).Where("recommendation_id = ?", recommendationID).Order("created_at desc").Find(&feedback).Error
	return feedback, err
}
