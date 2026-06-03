package httpctx

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize"`
	Offset   int `json:"offset"`
	Limit    int `json:"limit"`
}

func PaginationFromQuery(c *gin.Context) Pagination {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "pageSize", 20)
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return Pagination{
		Page:     page,
		PageSize: pageSize,
		Offset:   (page - 1) * pageSize,
		Limit:    pageSize,
	}
}

func queryInt(c *gin.Context, key string, fallback int) int {
	value := c.Query(key)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
}
