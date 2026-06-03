package workout

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

func (s Service) ListExercises(ctx context.Context) ([]models.Exercise, error) {
	return s.repo.ListExercises(ctx)
}

func (s Service) GetExercise(ctx context.Context, id string) (models.Exercise, error) {
	return s.repo.GetExercise(ctx, id)
}

func (s Service) CreateExercise(ctx context.Context, req ExerciseRequest) (models.Exercise, error) {
	exercise := models.Exercise{
		ID:             uuid.NewString(),
		Name:           req.Name,
		Category:       req.Category,
		MuscleGroup:    req.MuscleGroup,
		Equipment:      req.Equipment,
		IntensityLevel: req.IntensityLevel,
	}
	return exercise, s.repo.SaveExercise(ctx, &exercise)
}

func (s Service) UpdateExercise(ctx context.Context, id string, req ExerciseRequest) (models.Exercise, error) {
	exercise, err := s.repo.GetExercise(ctx, id)
	if err != nil {
		return exercise, err
	}
	exercise.Name = req.Name
	exercise.Category = req.Category
	exercise.MuscleGroup = req.MuscleGroup
	exercise.Equipment = req.Equipment
	exercise.IntensityLevel = req.IntensityLevel
	return exercise, s.repo.SaveExercise(ctx, &exercise)
}

func (s Service) DeleteExercise(ctx context.Context, id string) error {
	rows, err := s.repo.DeleteExercise(ctx, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s Service) ListSessions(ctx context.Context, userID string) ([]models.WorkoutSession, error) {
	return s.repo.ListSessions(ctx, userID)
}

func (s Service) GetSession(ctx context.Context, userID, id string) (models.WorkoutSession, error) {
	return s.repo.GetSession(ctx, userID, id)
}

func (s Service) CreateSession(ctx context.Context, userID string, req SessionRequest) (models.WorkoutSession, Indicators, error) {
	workoutDate, err := parseDate(req.WorkoutDate, "workoutDate")
	if err != nil {
		return models.WorkoutSession{}, Indicators{}, err
	}

	session := models.WorkoutSession{
		ID:          uuid.NewString(),
		UserID:      userID,
		WorkoutDate: workoutDate,
		Category:    req.Category,
		DurationMin: req.DurationMin,
		Notes:       req.Notes,
		Entries:     buildEntries("", req.Entries),
	}
	if err := s.repo.CreateSession(ctx, &session); err != nil {
		return session, Indicators{}, err
	}
	session, err = s.repo.GetSession(ctx, userID, session.ID)
	if err != nil {
		return session, Indicators{}, err
	}

	shared.InvalidateDashboard(ctx, s.redis, userID)
	return session, Analyze(session), nil
}

func (s Service) UpdateSession(ctx context.Context, userID, id string, req SessionRequest) (models.WorkoutSession, Indicators, error) {
	workoutDate, err := parseDate(req.WorkoutDate, "workoutDate")
	if err != nil {
		return models.WorkoutSession{}, Indicators{}, err
	}

	session, err := s.repo.GetSession(ctx, userID, id)
	if err != nil {
		return session, Indicators{}, err
	}
	session.WorkoutDate = workoutDate
	session.Category = req.Category
	session.DurationMin = req.DurationMin
	session.Notes = req.Notes

	entries := buildEntries(session.ID, req.Entries)
	if err := s.repo.ReplaceSession(ctx, &session, entries); err != nil {
		return session, Indicators{}, err
	}
	session, err = s.repo.GetSession(ctx, userID, id)
	if err != nil {
		return session, Indicators{}, err
	}

	shared.InvalidateDashboard(ctx, s.redis, userID)
	return session, Analyze(session), nil
}

func (s Service) DeleteSession(ctx context.Context, userID, id string) error {
	rows, err := s.repo.DeleteSession(ctx, userID, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return nil
}

func buildEntries(sessionID string, items []EntryRequest) []models.WorkoutEntry {
	entries := make([]models.WorkoutEntry, 0, len(items))
	for _, item := range items {
		entries = append(entries, models.WorkoutEntry{
			ID:         uuid.NewString(),
			SessionID:  sessionID,
			ExerciseID: item.ExerciseID,
			Sets:       item.Sets,
			Reps:       item.Reps,
			WeightKg:   item.WeightKg,
			RestSec:    item.RestSec,
		})
	}
	return entries
}

func parseDate(value, field string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, errors.New(field + " must be YYYY-MM-DD")
	}
	return parsed, nil
}
