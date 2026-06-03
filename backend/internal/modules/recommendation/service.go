package recommendation

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/analytics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Service struct {
	cfg      config.Config
	db       *gorm.DB
	repo     Repository
	provider Provider
	redis    *redis.Client
}

func NewService(cfg config.Config, db *gorm.DB, redisClient *redis.Client) Service {
	return Service{
		cfg:      cfg,
		db:       db,
		repo:     NewRepository(db),
		provider: NewProvider(cfg),
		redis:    redisClient,
	}
}

func (s Service) List(ctx context.Context, userID string) ([]models.AIRecommendation, error) {
	return s.repo.List(ctx, userID)
}

func (s Service) Get(ctx context.Context, userID, id string) (models.AIRecommendation, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s Service) Generate(ctx context.Context, userID string, req GenerateRequest) (models.AIRecommendation, error) {
	if err := s.checkRateLimit(ctx, userID); err != nil {
		return models.AIRecommendation{}, err
	}
	summary, err := analytics.BuildSummary(ctx, s.db, userID, "weekly")
	if err != nil {
		return models.AIRecommendation{}, err
	}

	contextPayload, _ := json.MarshalIndent(gin.H{
		"provider": s.cfg.AIProvider,
		"type":     req.Type,
		"summary":  summary,
	}, "", "  ")

	content, err := s.provider.Generate(ctx, req.Type, summary)
	if err != nil {
		return models.AIRecommendation{}, err
	}
	if !IsSafeGeneralWellness(content) {
		return models.AIRecommendation{}, errors.New("recommendation did not pass safety validation")
	}

	recommendation := models.AIRecommendation{
		ID:            uuid.NewString(),
		UserID:        userID,
		Type:          req.Type,
		PromptContext: string(contextPayload),
		Content:       content,
	}
	return recommendation, s.repo.Save(ctx, &recommendation)
}

func (s Service) checkRateLimit(ctx context.Context, userID string) error {
	if s.redis == nil || s.cfg.AIRateLimitPerHour <= 0 {
		return nil
	}
	key := "rate:ai:" + userID + ":" + strconv.FormatInt(time.Now().Unix()/3600, 10)
	count, err := s.redis.Incr(ctx, key).Result()
	if err != nil {
		return nil
	}
	if count == 1 {
		_ = s.redis.Expire(ctx, key, time.Hour).Err()
	}
	if count > int64(s.cfg.AIRateLimitPerHour) {
		return errors.New("AI recommendation rate limit exceeded")
	}
	return nil
}

func (s Service) Delete(ctx context.Context, userID, id string) error {
	rows, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s Service) AddFeedback(ctx context.Context, userID, recommendationID string, req FeedbackRequest) (models.RecommendationFeedback, error) {
	if _, err := s.repo.Get(ctx, userID, recommendationID); err != nil {
		return models.RecommendationFeedback{}, err
	}
	feedback := models.RecommendationFeedback{
		ID:               uuid.NewString(),
		RecommendationID: recommendationID,
		Rating:           req.Rating,
		Suitability:      req.Suitability,
		Comment:          req.Comment,
	}
	return feedback, s.repo.CreateFeedback(ctx, &feedback)
}

func (s Service) ListFeedback(ctx context.Context, userID, recommendationID string) ([]models.RecommendationFeedback, error) {
	if _, err := s.repo.Get(ctx, userID, recommendationID); err != nil {
		return nil, err
	}
	return s.repo.ListFeedback(ctx, recommendationID)
}
