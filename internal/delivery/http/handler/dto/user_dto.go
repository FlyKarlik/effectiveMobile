package http_dto

import (
	"strconv"

	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	"github.com/gin-gonic/gin"
)

func GetUserFilterFromQuery(c *gin.Context) domain.UserFilter {
	filter := domain.UserFilter{}

	if name := c.Query("name"); name != "" {
		filter.Name = &name
	}

	if surname := c.Query("surname"); surname != "" {
		filter.Surname = &surname
	}

	if patronymic := c.Query("patronymic"); patronymic != "" {
		filter.Patronymic = &patronymic
	}

	if ageStr := c.Query("age"); ageStr != "" {
		if age, err := strconv.Atoi(ageStr); err == nil {
			int64Age := int64(age)
			filter.Age = &int64Age
		}
	}

	if nationality := c.Query("nationality"); nationality != "" {
		filter.Nationality = &nationality
	}

	if sex := c.Query("sex"); sex != "" {
		filter.Sex = (*domain.SexEnum)(&sex)
	}

	return filter
}
