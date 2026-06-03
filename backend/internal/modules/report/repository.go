package report

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

func (r Repository) List(ctx context.Context, userID string) ([]models.ProgressReport, error) {
	var reports []models.ProgressReport
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("generated_at desc").Limit(50).Find(&reports).Error
	return reports, err
}

func (r Repository) Get(ctx context.Context, userID, id string) (models.ProgressReport, error) {
	var report models.ProgressReport
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&report).Error
	return report, err
}

func (r Repository) Save(ctx context.Context, report *models.ProgressReport) error {
	return r.db.WithContext(ctx).Save(report).Error
}

func (r Repository) Delete(ctx context.Context, userID, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&models.ProgressReport{}, "id = ? AND user_id = ?", id, userID)
	return result.RowsAffected, result.Error
}
