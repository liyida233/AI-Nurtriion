package httpctx

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const userIDKey = "userID"

func SetUserID(c *gin.Context, userID string) {
	c.Set(userIDKey, userID)
}

func UserID(c *gin.Context) (string, error) {
	value, exists := c.Get(userIDKey)
	if !exists {
		return "", errors.New("authenticated user missing from context")
	}
	userID, ok := value.(string)
	if !ok || userID == "" {
		return "", errors.New("authenticated user is invalid")
	}
	return userID, nil
}
