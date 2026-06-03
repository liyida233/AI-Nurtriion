package middleware

import (
	"net/http"

	"ai-nutrition/backend/internal/httpctx"
	"ai-nutrition/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func AdminRequired(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, err := httpctx.UserID(c)
		if err != nil {
			httpctx.Error(c, http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}

		var user models.User
		if err := db.Select("id", "role").First(&user, "id = ?", userID).Error; err != nil {
			httpctx.Error(c, http.StatusUnauthorized, "authenticated user not found")
			c.Abort()
			return
		}
		if user.Role != "admin" {
			httpctx.Error(c, http.StatusForbidden, "admin role required")
			c.Abort()
			return
		}
		c.Next()
	}
}
