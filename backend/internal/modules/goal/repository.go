package goal

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

func (r Repository) List(ctx context.Context, userID string) ([]models.Goal, error) {
	var goals []models.Goal
	err := r.db.WithContext(ctx).Preload("Milestones").Where("user_id = ?", userID).Order("created_at desc").Find(&goals).Error
	return goals, err
}

func (r Repository) Get(ctx context.Context, userID, id string) (models.Goal, error) {
	var goal models.Goal
	err := r.db.WithContext(ctx).Preload("Milestones").Where("id = ? AND user_id = ?", id, userID).First(&goal).Error
	return goal, err
}

func (r Repository) Save(ctx context.Context, goal *models.Goal) error {
	return r.db.WithContext(ctx).Save(goal).Error
}

func (r Repository) Create(ctx context.Context, goal *models.Goal) error {
	return r.db.WithContext(ctx).Create(goal).Error
}

func (r Repository) Delete(ctx context.Context, userID, id string) (int64, error) {
	var rows int64
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.GoalMilestone{}, "goal_id = ?", id).Error; err != nil {
			return err
		}
		result := tx.Delete(&models.Goal{}, "id = ? AND user_id = ?", id, userID)
		rows = result.RowsAffected
		return result.Error
	})
	return rows, err
}

func (r Repository) GetMilestone(ctx context.Context, goalID, milestoneID string) (models.GoalMilestone, error) {
	var milestone models.GoalMilestone
	err := r.db.WithContext(ctx).Where("id = ? AND goal_id = ?", milestoneID, goalID).First(&milestone).Error
	return milestone, err
}

func (r Repository) SaveMilestone(ctx context.Context, milestone *models.GoalMilestone) error {
	return r.db.WithContext(ctx).Save(milestone).Error
}
