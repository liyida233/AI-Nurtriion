package report

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"time"

	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/analytics"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GeneratedReport struct {
	Report  models.ProgressReport      `json:"report"`
	Metrics analytics.DashboardSummary `json:"metrics"`
}

type Service struct {
	db   *gorm.DB
	repo Repository
}

func NewService(db *gorm.DB) Service {
	return Service{db: db, repo: NewRepository(db)}
}

func (s Service) List(ctx context.Context, userID string) ([]models.ProgressReport, error) {
	return s.repo.List(ctx, userID)
}

func (s Service) Get(ctx context.Context, userID, id string) (models.ProgressReport, error) {
	return s.repo.Get(ctx, userID, id)
}

func (s Service) Generate(ctx context.Context, userID string, req GenerateRequest) (GeneratedReport, error) {
	start, end, period, err := analytics.ResolveRange(req.PeriodType, req.StartDate, req.EndDate)
	if err != nil {
		return GeneratedReport{}, err
	}
	summary, err := analytics.NewService(s.db, nil).BuildSummary(ctx, userID, period, start, end, false)
	if err != nil {
		return GeneratedReport{}, err
	}
	report := models.ProgressReport{
		ID:          uuid.NewString(),
		UserID:      userID,
		PeriodType:  period,
		GeneratedAt: time.Now(),
		Summary:     BuildNarrativeSummary(summary),
		FileURL:     "",
	}
	fileURL, err := writeHTMLReport(report.ID, summary, report.Summary)
	if err != nil {
		return GeneratedReport{}, err
	}
	report.FileURL = fileURL
	if err := s.repo.Save(ctx, &report); err != nil {
		return GeneratedReport{}, err
	}
	return GeneratedReport{Report: report, Metrics: summary}, nil
}

func (s Service) Delete(ctx context.Context, userID, id string) error {
	report, getErr := s.repo.Get(ctx, userID, id)
	rows, err := s.repo.Delete(ctx, userID, id)
	if err != nil {
		return err
	}
	if rows == 0 {
		return gorm.ErrRecordNotFound
	}
	if getErr == nil && report.FileURL != "" {
		_ = os.Remove(reportHTMLPath(report.ID))
	}
	return nil
}

func (s Service) ReportFilePath(ctx context.Context, userID, id string) (string, error) {
	report, err := s.repo.Get(ctx, userID, id)
	if err != nil {
		return "", err
	}
	if report.FileURL == "" {
		return "", errors.New("report file not available")
	}
	return reportHTMLPath(report.ID), nil
}

func writeHTMLReport(reportID string, summary analytics.DashboardSummary, narrative string) (string, error) {
	dir := filepath.Join("storage", "reports")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	filename := reportID + ".html"
	path := filepath.Join(dir, filename)
	if err := os.WriteFile(path, []byte(BuildHTML(reportID, summary, narrative)), 0644); err != nil {
		return "", err
	}
	return "/api/reports/" + reportID + "/download", nil
}

func reportHTMLPath(reportID string) string {
	return filepath.Join("storage", "reports", reportID+".html")
}
