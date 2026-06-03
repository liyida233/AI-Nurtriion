package analytics

import (
	"context"
	"time"

	"ai-nutrition/backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) WorkoutSessions(ctx context.Context, userID string, start time.Time) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := r.db.WithContext(ctx).Preload("Entries.Exercise").Where("user_id = ? AND workout_date >= ?", userID, start).Find(&sessions).Error
	return sessions, err
}

func (r Repository) Meals(ctx context.Context, userID string, start time.Time) ([]models.MealLog, error) {
	var meals []models.MealLog
	err := r.db.WithContext(ctx).Where("user_id = ? AND meal_time >= ?", userID, start).Find(&meals).Error
	return meals, err
}

func (r Repository) Profile(ctx context.Context, userID string) (models.UserProfile, error) {
	var profile models.UserProfile
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&profile).Error
	return profile, err
}

func (r Repository) BodyRecords(ctx context.Context, userID string) ([]models.BodyRecord, error) {
	var records []models.BodyRecord
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("record_date asc").Limit(30).Find(&records).Error
	return records, err
}

func (r Repository) ActiveGoalCount(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&models.Goal{}).Where("user_id = ? AND status = ?", userID, "active").Count(&count).Error
	return count, err
}

func (r Repository) Goals(ctx context.Context, userID string) ([]models.Goal, error) {
	var goals []models.Goal
	err := r.db.WithContext(ctx).Preload("Milestones").Where("user_id = ?", userID).Find(&goals).Error
	return goals, err
}

func (r Repository) SaveSnapshot(ctx context.Context, snapshot *models.AnalyticsSnapshot) error {
	return r.db.WithContext(ctx).Create(snapshot).Error
}

func (r Repository) ListSnapshots(ctx context.Context, userID string) ([]models.AnalyticsSnapshot, error) {
	var snapshots []models.AnalyticsSnapshot
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at desc").Limit(50).Find(&snapshots).Error
	return snapshots, err
}
