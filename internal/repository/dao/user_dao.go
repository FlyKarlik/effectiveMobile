package dao

import (
	"database/sql"
	"time"

	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	"github.com/FlyKarlik/effectiveMobile/pkg/database/postgres"
	"github.com/google/uuid"
)

type UserDAO struct {
	ID          uuid.UUID
	Name        string
	Surname     string
	Patronymic  sql.NullString
	Nationality sql.NullString
	Age         sql.NullInt64
	Sex         sql.NullString
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (u *UserDAO) ToDomain() domain.User {
	return domain.User{
		ID:          u.ID,
		CreateAt:    u.CreatedAt,
		UpdateAt:    u.UpdatedAt,
		Name:        u.Name,
		Surname:     u.Surname,
		Patronymic:  postgres.FromNullString(u.Patronymic),
		Nationality: postgres.FromNullString(u.Nationality),
		Age:         postgres.FromNullInt64(u.Age),
		Sex:         (*domain.SexEnum)(postgres.FromNullString(u.Sex)),
	}
}

type CreateUserInputDAO struct {
	Name        string
	Surname     string
	Patronymic  sql.NullString
	Nationality sql.NullString
	Age         sql.NullInt64
	Sex         sql.NullString
}

func (c *CreateUserInputDAO) FromDomain(domain domain.CreateUserInput) {
	c.Name = domain.Name
	c.Surname = domain.Surname
	c.Patronymic = postgres.ToNullString(domain.Patronymic)
	c.Nationality = postgres.ToNullString(domain.Nationality)
	c.Age = postgres.ToNullInt64(domain.Age)
	c.Sex = postgres.ToNullString((*string)(domain.Sex))
}

type UpdateUserInputDAO struct {
	Name        sql.NullString
	Surname     sql.NullString
	Patronymic  sql.NullString
	Nationality sql.NullString
	Age         sql.NullInt64
	Sex         sql.NullString
}

func (u *UpdateUserInputDAO) FromDomain(domain domain.UpdateUserInput) {
	u.Name = postgres.ToNullString(domain.Name)
	u.Surname = postgres.ToNullString(domain.Surname)
	u.Patronymic = postgres.ToNullString(domain.Patronymic)
	u.Nationality = postgres.ToNullString(domain.Nationality)
	u.Age = postgres.ToNullInt64(domain.Age)
	u.Sex = postgres.ToNullString((*string)(domain.Sex))
}

type UserFilterDAO struct {
	Name        sql.NullString
	Surname     sql.NullString
	Patronymic  sql.NullString
	Nationality sql.NullString
	Sex         sql.NullString
	Age         sql.NullInt64
}

func (u *UserFilterDAO) FromDomain(domain domain.UserFilter) {
	u.Name = postgres.ToNullString(domain.Name)
	u.Surname = postgres.ToNullString(domain.Surname)
	u.Patronymic = postgres.ToNullString(domain.Patronymic)
	u.Nationality = postgres.ToNullString(domain.Nationality)
	u.Sex = postgres.ToNullString((*string)(domain.Sex))
	u.Age = postgres.ToNullInt64(domain.Age)
}
