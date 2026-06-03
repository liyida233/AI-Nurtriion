package admin

import (
	"context"

	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{db: db}
}

func (r Repository) ListUsers(ctx context.Context, pagination httpctx.Pagination) ([]models.User, error) {
	var users []models.User
	err := r.db.WithContext(ctx).
		Preload("Profile").
		Order("created_at desc").
		Offset(pagination.Offset).
		Limit(pagination.Limit).
		Find(&users).Error
	return users, err
}

func (r Repository) GetUser(ctx context.Context, id string) (models.User, error) {
	var user models.User
	err := r.db.WithContext(ctx).Preload("Profile").First(&user, "id = ?", id).Error
	return user, err
}

func (r Repository) SaveUser(ctx context.Context, user *models.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}
