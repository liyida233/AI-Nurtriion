package goal

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

func (s Service) List(ctx context.Context, userID string) ([]models.Goal, error) {
	return s.repo.List(ctx, userID)
}

func (s Service) Get(ctx context.Context, userID, id string) (models.Goal, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s Service) Create(ctx context.Context, userID string, req GoalRequest) (models.Goal, error) {
	deadline, err := parseDate(req.Deadline)
	if err != nil {
		return models.Goal{}, err
	}
	if err := FeasibilityWarning(req.GoalType, req.TargetMetric, req.TargetValue, deadline); err != nil {
		return models.Goal{}, err
	}
	goal := mapGoal(models.Goal{
		ID:     uuid.NewString(),
		UserID: userID,
		Status: "active",
	}, req, deadline)
	milestones, err := resolveMilestones(goal, req)
	if err != nil {
		return models.Goal{}, err
	}
	goal.Milestones = milestones
	if err := s.repo.Create(ctx, &goal); err != nil {
		return goal, err
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return s.repo.Get(ctx, userID, goal.ID)
}

func (s Service) Update(ctx context.Context, userID, id string, req GoalRequest) (models.Goal, error) {
	goal, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return goal, err
	}
	deadline, err := parseDate(req.Deadline)
	if err != nil {
		return goal, err
	}
	if err := FeasibilityWarning(req.GoalType, req.TargetMetric, req.TargetValue, deadline); err != nil {
		return goal, err
	}
	goal = mapGoal(goal, req, deadline)
	if err := s.repo.Save(ctx, &goal); err != nil {
		return goal, err
	}
	if req.Milestones != nil {
		milestones, err := resolveMilestones(goal, req)
		if err != nil {
			return goal, err
		}
		if err := s.repo.ReplaceMilestones(ctx, goal.ID, milestones); err != nil {
			return goal, err
		}
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return s.repo.Get(ctx, userID, id)
}

func (s Service) UpdateStatus(ctx context.Context, userID, id, status string) (models.Goal, error) {
	if status != "active" && status != "completed" && status != "paused" && status != "cancelled" {
		return models.Goal{}, errors.New("status must be active, completed, paused, or cancelled")
	}
	goal, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return goal, err
	}
	goal.Status = status
	if err := s.repo.Save(ctx, &goal); err != nil {
		return goal, err
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return s.repo.Get(ctx, userID, id)
}

func (s Service) UpdateMilestone(ctx context.Context, userID, goalID, milestoneID string, completed bool) (models.GoalMilestone, error) {
	if _, err := s.repo.Get(ctx, userID, goalID); err != nil {
		return models.GoalMilestone{}, err
	}
	milestone, err := s.repo.GetMilestone(ctx, goalID, milestoneID)
	if err != nil {
		return milestone, err
	}
	milestone.Completed = completed
	if err := s.repo.SaveMilestone(ctx, &milestone); err != nil {
		return milestone, err
	}
	shared.InvalidateDashboard(ctx, s.redis, userID)
	return milestone, nil
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

func mapGoal(goal models.Goal, req GoalRequest, deadline time.Time) models.Goal {
	goal.GoalType = req.GoalType
	goal.TargetMetric = req.TargetMetric
	goal.TargetValue = req.TargetValue
	goal.Deadline = deadline
	goal.Priority = req.Priority
	return goal
}

func parseDate(value string) (time.Time, error) {
	parsed, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, errors.New("deadline must be YYYY-MM-DD")
	}
	return parsed, nil
}

func buildMilestones(goal models.Goal) []models.GoalMilestone {
	now := time.Now()
	totalDays := goal.Deadline.Sub(now).Hours() / 24
	if totalDays <= 0 {
		return nil
	}
	milestoneCount := 3
	milestones := make([]models.GoalMilestone, 0, milestoneCount)
	for i := 1; i <= milestoneCount; i++ {
		ratio := float64(i) / float64(milestoneCount)
		milestones = append(milestones, models.GoalMilestone{
			ID:          uuid.NewString(),
			Title:       "Milestone " + string(rune('0'+i)),
			TargetValue: goal.TargetValue * ratio,
			DueDate:     now.Add(time.Duration(totalDays*ratio) * 24 * time.Hour),
			Completed:   false,
		})
	}
	return milestones
}

func resolveMilestones(goal models.Goal, req GoalRequest) ([]models.GoalMilestone, error) {
	if len(req.Milestones) == 0 {
		return buildMilestones(goal), nil
	}
	milestones := make([]models.GoalMilestone, 0, len(req.Milestones))
	for index, item := range req.Milestones {
		if item.Title == "" {
			item.Title = "Milestone " + string(rune('1'+index))
		}
		dueDate, err := parseDate(item.DueDate)
		if err != nil {
			return nil, errors.New("milestone dueDate must be YYYY-MM-DD")
		}
		milestones = append(milestones, models.GoalMilestone{
			ID:          uuid.NewString(),
			GoalID:      goal.ID,
			Title:       item.Title,
			TargetValue: item.TargetValue,
			DueDate:     dueDate,
			Completed:   false,
		})
	}
	return milestones, nil
}
