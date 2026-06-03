package middleware

import (
	"net/http"
	"strings"

	"ai-nutrition/backend/internal/config"
	"ai-nutrition/backend/internal/httpctx"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
)

func AuthRequired(cfg config.Config, redisClient *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			httpctx.Error(c, http.StatusUnauthorized, "missing bearer token")
			c.Abort()
			return
		}

		tokenValue := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (any, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			httpctx.Error(c, http.StatusUnauthorized, "invalid token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			httpctx.Error(c, http.StatusUnauthorized, "invalid token claims")
			c.Abort()
			return
		}

		sub, ok := claims["sub"].(string)
		if !ok || sub == "" {
			httpctx.Error(c, http.StatusUnauthorized, "invalid token subject")
			c.Abort()
			return
		}

		if redisClient != nil {
			if err := redisClient.Get(c.Request.Context(), "session:"+sub).Err(); err != nil {
				httpctx.Error(c, http.StatusUnauthorized, "session expired or logged out")
				c.Abort()
				return
			}
		}

		httpctx.SetUserID(c, sub)
		c.Next()
	}
}
