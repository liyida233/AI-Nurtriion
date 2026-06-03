package body

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

func (s Service) List(ctx context.Context, userID string) ([]models.BodyRecord, error) {
	return s.repo.List(ctx, userID)
}

func (s Service) Get(ctx context.Context, userID, id string) (models.BodyRecord, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s Service) Create(ctx context.Context, userID string, req RecordRequest) (models.BodyRecord, error) {
	recordDate, err := parseDate(req.RecordDate)
	if err != nil {
		return models.BodyRecord{}, err
	}
	record := models.BodyRecord{
		ID:         uuid.NewString(),
		UserID:     userID,
		RecordDate: recordDate,
		WeightKg:   req.WeightKg,
		Note:       req.Note,
	}
	if err := s.repo.Save(ctx, &record); err != nil {
		return record, err
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return record, nil
}

func (s Service) Update(ctx context.Context, userID, id string, req RecordRequest) (models.BodyRecord, error) {
	record, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return record, err
	}
	recordDate, err := parseDate(req.RecordDate)
	if err != nil {
		return record, err
	}
	record.RecordDate = recordDate
	record.WeightKg = req.WeightKg
	record.Note = req.Note
	if err := s.repo.Save(ctx, &record); err != nil {
		return record, err
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return record, nil
}

func (s Service) Delete(ctx context.Context, userID, id string) error {
	rows, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return nil
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, errors.New("recordDate must be YYYY-MM-DD")
	}
	return parsed, nil
}
