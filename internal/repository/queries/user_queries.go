package queries

import (
	"time"

	"github.com/FlyKarlik/effectiveMobile/internal/repository/dao"
	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
)

func BuildCountUsersQuery(filter dao.UserFilterDAO) (string, []interface{}, error) {
	builder := sq.Select("COUNT(*)").From(`"user"`).PlaceholderFormat(sq.Dollar)

	if filter.Name.Valid {
		builder = builder.Where(sq.ILike{"name": "%" + filter.Name.String + "%"})
	}

	if filter.Surname.Valid {
		builder = builder.Where(sq.ILike{"surname": "%" + filter.Surname.String + "%"})
	}

	if filter.Patronymic.Valid {
		builder = builder.Where(sq.ILike{"patronymic": "%" + filter.Patronymic.String + "%"})
	}

	if filter.Nationality.Valid {
		builder = builder.Where(sq.Eq{"nationality": filter.Nationality.String})
	}

	if filter.Sex.Valid {
		builder = builder.Where(sq.Eq{"sex": filter.Sex.String})
	}

	if filter.Age.Valid {
		builder = builder.Where(sq.Eq{"age": filter.Age.Int64})
	}

	return builder.ToSql()
}
func BuildSearchUsersQuery(filter dao.UserFilterDAO, pagination dao.PaginationDAO) (string, []interface{}, error) {
	builder := sq.Select("*").From(`"user"`).PlaceholderFormat(sq.Dollar)

	if filter.Name.Valid {
		builder = builder.Where(sq.ILike{"name": "%" + filter.Name.String + "%"})
	}

	if filter.Surname.Valid {
		builder = builder.Where(sq.ILike{"surname": "%" + filter.Surname.String + "%"})
	}

	if filter.Patronymic.Valid {
		builder = builder.Where(sq.ILike{"patronymic": "%" + filter.Patronymic.String + "%"})
	}

	if filter.Nationality.Valid {
		builder = builder.Where(sq.Eq{"nationality": filter.Nationality.String})
	}

	if filter.Sex.Valid {
		builder = builder.Where(sq.Eq{"sex": filter.Sex.String})
	}

	if filter.Age.Valid {
		builder = builder.Where(sq.Eq{"age": filter.Age.Int64})
	}

	if pagination.Limit.Valid {
		builder = builder.Limit(uint64(pagination.Limit.Int64))
	}
	if pagination.Offset.Valid {
		builder = builder.Offset(uint64(pagination.Offset.Int64))
	}
	return builder.ToSql()
}

func BuildCreateUserQuery(user dao.CreateUserInputDAO) (string, []interface{}, error) {
	values := map[string]interface{}{
		"name":    user.Name,
		"surname": user.Surname,
	}

	if user.Patronymic.Valid {
		values["patronymic"] = user.Patronymic.String
	}

	if user.Nationality.Valid {
		values["nationality"] = user.Nationality.String
	}

	if user.Age.Valid {
		values["age"] = user.Age.Int64
	}

	if user.Sex.Valid {
		values["sex"] = user.Sex.String
	}

	builder := sq.Insert(`"user"`).
		SetMap(values).
		PlaceholderFormat(sq.Dollar).
		Suffix("RETURNING *")

	return builder.ToSql()
}

func BuildUpdateUserQuery(id uuid.UUID, user dao.UpdateUserInputDAO) (string, []interface{}, error) {
	builder := sq.Update(`"user"`).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id})

	builder = builder.Set("updated_at", time.Now())

	if user.Name.Valid {
		builder = builder.Set("name", user.Name.String)
	}

	if user.Surname.Valid {
		builder = builder.Set("surname", user.Surname.String)
	}

	if user.Patronymic.Valid {
		builder = builder.Set("patronymic", user.Patronymic.String)
	}

	if user.Nationality.Valid {
		builder = builder.Set("nationality", user.Nationality.String)
	}

	if user.Age.Valid {
		builder = builder.Set("age", user.Age.Int64)
	}

	if user.Sex.Valid {
		builder = builder.Set("sex", user.Sex.String)
	}

	builder = builder.Suffix("RETURNING *")

	return builder.ToSql()
}

func BuildDeleteUserQuery(id uuid.UUID) (string, []interface{}, error) {
	builder := sq.Delete(`"user"`).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": id}).
		Suffix("RETURNING *")

	return builder.ToSql()
}
