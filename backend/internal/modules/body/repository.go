package body

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

func (r Repository) List(ctx context.Context, userID string) ([]models.BodyRecord, error) {
	var records []models.BodyRecord
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("record_date desc").Limit(100).Find(&records).Error
	return records, err
}

func (r Repository) Get(ctx context.Context, userID, id string) (models.BodyRecord, error) {
	var record models.BodyRecord
	err := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).First(&record).Error
	return record, err
}

func (r Repository) Save(ctx context.Context, record *models.BodyRecord) error {
	return r.db.WithContext(ctx).Save(record).Error
}

func (r Repository) Delete(ctx context.Context, userID, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&models.BodyRecord{}, "id = ? AND user_id = ?", id, userID)
	return result.RowsAffected, result.Error
}
