package body

import (
	"net/http"
	"time"

	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"ai-nutrition/backend/internal/modules/shared"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	db    *gorm.DB
	redis *redis.Client
}

type createBodyRecordRequest struct {
	RecordDate string  `json:"recordDate" binding:"required"`
	WeightKg   float64 `json:"weightKg" binding:"required,min=25,max=350"`
	Note       string  `json:"note"`
}

func RegisterRoutes(router *gin.RouterGroup, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{db: db, redis: redisClient}
	router.GET("", handler.List)
	router.POST("", handler.Create)
}

func (h Handler) List(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var records []models.BodyRecord
	if err := h.db.Where("user_id = ?", userID).Order("record_date desc").Limit(100).Find(&records).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not list body records")
		return
	}
	httpctx.OK(c, records)
}

func (h Handler) Create(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var req createBodyRecordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	recordDate, err := time.Parse("2006-01-02", req.RecordDate)
	if err != nil {
		httpctx.Error(c, http.StatusBadRequest, "recordDate must be YYYY-MM-DD")
		return
	}

	record := models.BodyRecord{
		ID:         uuid.NewString(),
		UserID:     userID,
		RecordDate: recordDate,
		WeightKg:   req.WeightKg,
		Note:       req.Note,
	}

	if err := h.db.Create(&record).Error; err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not save body record")
		return
	}

	shared.InvalidateDashboard(c.Request.Context(), h.redis, userID)
	httpctx.Created(c, record)
}
