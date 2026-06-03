package auth

import (
	"context"
	"net/http"
	"time"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Handler struct {
	cfg   config.Config
	db    *gorm.DB
	redis *redis.Client
}

type registerRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func RegisterRoutes(router *gin.RouterGroup, cfg config.Config, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{cfg: cfg, db: db, redis: redisClient}
	router.POST("/register", handler.Register)
	router.POST("/login", handler.Login)
}

func RegisterProtectedRoutes(router *gin.RouterGroup, cfg config.Config, db *gorm.DB, redisClient *redis.Client) {
	handler := Handler{cfg: cfg, db: db, redis: redisClient}
	router.GET("/me", handler.Me)
	router.POST("/logout", handler.Logout)
}

func (h Handler) Register(c *gin.Context) {
	var req registerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not secure password")
		return
	}

	user := models.User{
		ID:           uuid.NewString(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: string(hash),
		Role:         "user",
	}

	if err := h.db.Create(&user).Error; err != nil {
		httpctx.Error(c, http.StatusConflict, "email already exists")
		return
	}

	token, err := h.issueToken(c.Request.Context(), user.ID)
	if err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not issue token")
		return
	}

	httpctx.Created(c, gin.H{"token": token, "user": user})
}

func (h Handler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		httpctx.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var user models.User
	if err := h.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		httpctx.Error(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		httpctx.Error(c, http.StatusUnauthorized, "invalid email or password")
		return
	}

	token, err := h.issueToken(c.Request.Context(), user.ID)
	if err != nil {
		httpctx.Error(c, http.StatusInternalServerError, "could not issue token")
		return
	}

	httpctx.OK(c, gin.H{"token": token, "user": user})
}

func (h Handler) Me(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}

	var user models.User
	if err := h.db.Preload("Profile").First(&user, "id = ?", userID).Error; err != nil {
		httpctx.Error(c, http.StatusNotFound, "user not found")
		return
	}
	httpctx.OK(c, user)
}

func (h Handler) Logout(c *gin.Context) {
	userID, err := httpctx.UserID(c)
	if err != nil {
		httpctx.Error(c, http.StatusUnauthorized, err.Error())
		return
	}
	if h.redis != nil {
		_ = h.redis.Del(c.Request.Context(), "session:"+userID).Err()
	}
	httpctx.OK(c, gin.H{"loggedOut": true})
}

func (h Handler) issueToken(ctx context.Context, userID string) (string, error) {
	expiresAt := time.Now().Add(24 * time.Hour)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": userID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	})

	signed, err := token.SignedString([]byte(h.cfg.JWTSecret))
	if err != nil {
		return "", err
	}

	if h.redis != nil {
		_ = h.redis.Set(ctx, "session:"+userID, "active", time.Until(expiresAt)).Err()
	}

	return signed, nil
}
