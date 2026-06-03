package recommendation

import (
	"encoding/json"
	"net/http"
	"strings"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/analytics"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	cfg   config.Config
	db    *gorm.DB
	redis *redis.Client
}

type generateRequest struct {
	Type string `json:"type" binding:"required"`
}

type feedbackRequest struct {
	Rating      string `json:"rating" binding:"required"`
	Suitability string `json:"suitability"`
	Comment     string `json:"comment"`
}

func RegisterRoutes(router *gin.RouterGroup, cfg config.Config, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{cfg: cfg, db: db, redis: redisClient}
	router.GET("", handler.List)
	router.POST("/generate", handler.Generate)
	router.POST("/:id/feedback", handler.Feedback)
}

func (h Handler) List(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var recommendations []models.AIRecommendation
	if err := h.db.Where("user_id = ?", userID).Order("created_at desc").Limit(20).Find(&recommendations).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not list recommendations")
		return
	}
	httpctx.OK(c, recommendations)
}

func (h Handler) Generate(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var req generateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	summary, err := analytics.BuildSummary(c.Request.Context(), h.db, userID, "weekly")
	if err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not prepare recommendation context")
		return
	}

	contextPayload, _ := json.MarshalIndent(gin.H{
		"provider": h.cfg.AIProvider,
		"type":     req.Type,
		"summary":  summary,
	}, "", "  ")

	content := mockRecommendation(req.Type, summary)
	if !isSafeGeneralWellness(content) {
		httpctx.Error(c, http.StatusUnprocessableEntity, "recommendation did not pass safety validation")
		return
	}

	recommendation := models.AIRecommendation{
		ID:            uuid.NewString(),
		UserID:        userID,
		Type:          req.Type,
		PromptContext: string(contextPayload),
		Content:       content,
	}

	if err := h.db.Create(&recommendation).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not save recommendation")
		return
	}

	httpctx.Created(c, recommendation)
}

func (h Handler) Feedback(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	recommendationID := c.Param("id")
	var recommendation models.AIRecommendation
	if err := h.db.Where("id = ? AND user_id = ?", recommendationID, userID).First(&recommendation).Error; err != nil {
		httpctx.Error(c, http.StatusNotFound, "recommendation not found")
		return
	}

	var req feedbackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	feedback := models.RecommendationFeedback{
		ID:               uuid.NewString(),
		RecommendationID: recommendationID,
		Rating:           req.Rating,
		Suitability:      req.Suitability,
		Comment:          req.Comment,
	}
	if err := h.db.Create(&feedback).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not save feedback")
		return
	}

	httpctx.Created(c, feedback)
}

func mockRecommendation(kind string, summary analytics.DashboardSummary) string {
	lines := []string{
		"Weekly recommendation:",
		"- Keep this as general wellness guidance, not medical advice.",
	}

	if kind == "meal" || kind == "weekly" {
		if summary.Protein < 420 {
			lines = append(lines, "- Protein intake looks low for the week; add lean protein to one or two meals per day.")
		}
		if summary.CalorieBalance > 2500 {
			lines = append(lines, "- Recent calorie balance is high; consider smaller portions or lower-calorie snacks.")
		}
	}

	if kind == "workout" || kind == "weekly" {
		if summary.WorkoutSessions < 3 {
			lines = append(lines, "- Workout consistency can improve; schedule three short sessions before increasing intensity.")
		} else {
			lines = append(lines, "- Workout frequency is on track; progress by adding small increases in reps, sets, or weight.")
		}
	}

	if summary.WeightTrend == "increasing" && summary.CalorieBalance > 0 {
		lines = append(lines, "- Body trend and calorie balance both point upward; check whether this matches your current goal.")
	}
	if summary.ActiveGoals == 0 {
		lines = append(lines, "- Add one SMART goal so future recommendations can be more targeted.")
	}

	return strings.Join(lines, "\n")
}

func isSafeGeneralWellness(content string) bool {
	blocked := []string{"diagnose", "cure", "starve", "extreme fasting"}
	lowered := strings.ToLower(content)
	for _, word := range blocked {
		if strings.Contains(lowered, word) {
			return false
		}
	}
	return true
}
