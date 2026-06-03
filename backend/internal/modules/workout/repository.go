package workout

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

func (r Repository) ListExercises(ctx context.Context) ([]models.Exercise, error) {
	var exercises []models.Exercise
	err := r.db.WithContext(ctx).Order("name asc").Find(&exercises).Error
	return exercises, err
}

func (r Repository) GetExercise(ctx context.Context, id string) (models.Exercise, error) {
	var exercise models.Exercise
	err := r.db.WithContext(ctx).First(&exercise, "id = ?", id).Error
	return exercise, err
}

func (r Repository) SaveExercise(ctx context.Context, exercise *models.Exercise) error {
	return r.db.WithContext(ctx).Save(exercise).Error
}

func (r Repository) DeleteExercise(ctx context.Context, id string) (int64, error) {
	result := r.db.WithContext(ctx).Delete(&models.Exercise{}, "id = ?", id)
	return result.RowsAffected, result.Error
}

func (r Repository) ListSessions(ctx context.Context, userID string) ([]models.WorkoutSession, error) {
	var sessions []models.WorkoutSession
	err := r.db.WithContext(ctx).
		Preload("Entries.Exercise").
		Where("user_id = ?", userID).
		Order("workout_date desc").
		Limit(100).
		Find(&sessions).Error
	return sessions, err
}

func (r Repository) GetSession(ctx context.Context, userID, id string) (models.WorkoutSession, error) {
	var session models.WorkoutSession
	err := r.db.WithContext(ctx).
		Preload("Entries.Exercise").
		Where("id = ? AND user_id = ?", id, userID).
		First(&session).Error
	return session, err
}

func (r Repository) CreateSession(ctx context.Context, session *models.WorkoutSession) error {
	return r.db.WithContext(ctx).Create(session).Error
}

func (r Repository) ReplaceSession(ctx context.Context, session *models.WorkoutSession, entries []models.WorkoutEntry) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(session).Error; err != nil {
			return err
		}
		if err := tx.Delete(&models.WorkoutEntry{}, "session_id = ?", session.ID).Error; err != nil {
			return err
		}
		return tx.Create(&entries).Error
	})
}

func (r Repository) DeleteSession(ctx context.Context, userID, id string) (int64, error) {
	var rows int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.WorkoutEntry{}, "session_id = ?", id).Error; err != nil {
			return err
		}
		result := tx.Delete(&models.WorkoutSession{}, "id = ? AND user_id = ?", id, userID)
		rows = result.RowsAffected
		return result.Error
	})
	return rows, err
}
