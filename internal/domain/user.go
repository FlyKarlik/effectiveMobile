package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	CreateAt    time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Surname     string    `json:"surname"`
	Nationality *string   `json:"nationality,omitempty"`
	Patronymic  *string   `json:"patronymic,omitempty"`
	Sex         *SexEnum  `json:"sex,omitempty"`
	Age         *int64    `json:"age,omitempty"`
}

type CreateUserInput struct {
	Name        string   `json:"name" binding:"required"`
	Surname     string   `json:"surname" binding:"required"`
	Patronymic  *string  `json:"patronymic,omitempty"`
	Nationality *string  `json:"-"`
	Age         *int64   `json:"-"`
	Sex         *SexEnum `json:"-"`
}

type UpdateUserInput struct {
	Name        *string  `json:"name,omitempty"`
	Surname     *string  `json:"surname,omitempty"`
	Patronymic  *string  `json:"patronymic,omitempty"`
	Nationality *string  `json:"nationality,omitempty"`
	Age         *int64   `json:"age,omitempty"`
	Sex         *SexEnum `json:"sex,omitempty"`
}

type UserFilter struct {
	Name        *string
	Surname     *string
	Patronymic  *string
	Nationality *string
	Sex         *SexEnum
	Age         *int64
}
