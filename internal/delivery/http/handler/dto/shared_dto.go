package http_dto

import (
	"strconv"

	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	"github.com/gin-gonic/gin"
)

func validateLimit(limit int64) int64 {
	if limit <= 0 || limit > 100 {
		return 10
	}
	return limit
}

func GetPaginationFromQuery(c *gin.Context) domain.Pagination {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.ParseInt(limitStr, 10, 64)
	if err != nil {
		limit = 10
	}
	limit = validateLimit(limit)

	offset, err := strconv.ParseInt(offsetStr, 10, 64)
	if err != nil || offset < 0 {
		offset = 0
	}

	return domain.Pagination{
		Limit:  limit,
		Offset: offset,
	}
}
