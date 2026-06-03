package admin

import (
	"context"
	"errors"

	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"gorm.io/gorm"
)

type Service struct {
	repo Repository
}

func NewService(db *gorm.DB) Service {
	return Service{repo: NewRepository(db)}
}

func (s Service) ListUsers(ctx context.Context, pagination httpctx.Pagination) ([]models.User, error) {
	return s.repo.ListUsers(ctx, pagination)
}

func (s Service) GetUser(ctx context.Context, id string) (models.User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s Service) UpdateRole(ctx context.Context, id string, role string) (models.User, error) {
	if role != "user" && role != "admin" {
		return models.User{}, errors.New("role must be user or admin")
	}
	user, err := s.repo.GetUser(ctx, id)
	if err != nil {
		return user, err
	}
	user.Role = role
	if err := s.repo.SaveUser(ctx, &user); err != nil {
		return user, err
	}
	return s.repo.GetUser(ctx, id)
}
