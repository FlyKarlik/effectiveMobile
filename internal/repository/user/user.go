package user_repo

import (
	"context"

	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	"github.com/FlyKarlik/effectiveMobile/internal/repository/dao"
	"github.com/FlyKarlik/effectiveMobile/internal/repository/queries"
	"github.com/FlyKarlik/effectiveMobile/pkg/database/postgres"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
	"github.com/google/uuid"
)

type IUserRepository interface {
	CountUsers(ctx context.Context, filter domain.UserFilter) (int64, error)
	SearchUsers(ctx context.Context, pagination domain.Pagination, filter domain.UserFilter) ([]domain.User, error)
	CreateUser(ctx context.Context, input domain.CreateUserInput) (domain.User, error)
	UpdateUserByID(ctx context.Context, ID uuid.UUID, input domain.UpdateUserInput) (domain.User, error)
	DeleteUserByID(ctx context.Context, ID uuid.UUID) (domain.User, error)
}

type userRepo struct {
	logger logger.Logger
	q      postgres.Querier
}

func New(logger logger.Logger, q postgres.Querier) IUserRepository {
	return &userRepo{
		logger: logger,
		q:      q,
	}
}

func (u *userRepo) CountUsers(ctx context.Context, filter domain.UserFilter) (int64, error) {
	const layer string = "repository"
	const method = "CountUsers"

	u.logger.Debug(layer, method, "started", "filter", filter)

	filterDAO := new(dao.UserFilterDAO)
	filterDAO.FromDomain(filter)

	query, args, err := queries.BuildCountUsersQuery(*filterDAO)
	if err != nil {
		u.logger.Error(layer, method, "failed to build query", err, "filter", filter)
		return 0, err
	}

	u.logger.Debug(layer, method, "query built", "query", query, "args", args)

	var count int64
	err = u.q.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		u.logger.Error(layer, method, "query execution failed", err, "query", query, "args", args)
		return 0, err
	}

	u.logger.Debug(layer, method, "successfully completed", "count", count)
	return count, nil
}

func (u *userRepo) SearchUsers(ctx context.Context, pagination domain.Pagination, filter domain.UserFilter) ([]domain.User, error) {
	const method = "SearchUsers"
	const layer string = "repository"

	u.logger.Debug(layer, method, "started", "pagination", pagination, "filter", filter)

	paginationDAO := new(dao.PaginationDAO)
	filterDAO := new(dao.UserFilterDAO)

	paginationDAO.FromDomain(pagination)
	filterDAO.FromDomain(filter)

	query, args, err := queries.BuildSearchUsersQuery(*filterDAO, *paginationDAO)
	if err != nil {
		u.logger.Error(layer, method, "failed to build query", err, "filter", filter, "pagination", pagination)
		return nil, err
	}

	u.logger.Debug(layer, method, "query built", "query", query, "args", args)

	rows, err := u.q.Query(ctx, query, args...)
	if err != nil {
		u.logger.Error(layer, method, "query execution failed", err, "query", query, "args", args)
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		var user dao.UserDAO
		err := rows.Scan(
			&user.ID,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.Name,
			&user.Surname,
			&user.Nationality,
			&user.Patronymic,
			&user.Sex,
			&user.Age,
		)
		if err != nil {
			u.logger.Error(layer, method, "row scan failed", err, "query", query)
			return nil, err
		}
		users = append(users, user.ToDomain())
	}

	if err := rows.Err(); err != nil {
		u.logger.Error(layer, method, "rows iteration error", err)
		return nil, err
	}

	u.logger.Debug(layer, method, "successfully completed", "users_count", len(users))
	return users, nil
}

func (u *userRepo) CreateUser(ctx context.Context, input domain.CreateUserInput) (domain.User, error) {
	const layer string = "repository"
	const method string = "CreateUser"

	u.logger.Debug(layer, method, "started", "input", input)

	createUserInputDAO := new(dao.CreateUserInputDAO)
	createUserInputDAO.FromDomain(input)

	query, args, err := queries.BuildCreateUserQuery(*createUserInputDAO)
	if err != nil {
		u.logger.Error(layer, method, "failed to build query", err, "input", input)
		return domain.User{}, err
	}

	u.logger.Debug(layer, method, "query built", "query", query, "args", args)

	var user dao.UserDAO
	err = u.q.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Surname,
		&user.Nationality,
		&user.Patronymic,
		&user.Sex,
		&user.Age,
	)
	if err != nil {
		u.logger.Error(layer, method, "query execution failed", err, "query", query, "args", args)
		return domain.User{}, err
	}

	result := user.ToDomain()
	u.logger.Debug(layer, method, "successfully completed", "created_user", result)
	return result, nil
}

func (u *userRepo) UpdateUserByID(ctx context.Context, ID uuid.UUID, input domain.UpdateUserInput) (domain.User, error) {
	const layer string = "repository"
	const method = "UpdateUserByID"

	u.logger.Debug(layer, method, "started", "id", ID, "input", input)

	updateUserInputDAO := new(dao.UpdateUserInputDAO)
	updateUserInputDAO.FromDomain(input)

	query, args, err := queries.BuildUpdateUserQuery(ID, *updateUserInputDAO)
	if err != nil {
		u.logger.Error(layer, method, "failed to build query", err, "id", ID, "input", input)
		return domain.User{}, err
	}

	u.logger.Debug(layer, method, "query built", "query", query, "args", args)

	var user dao.UserDAO
	err = u.q.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Surname,
		&user.Nationality,
		&user.Patronymic,
		&user.Sex,
		&user.Age,
	)
	if err != nil {
		u.logger.Error(layer, method, "query execution failed", err, "query", query, "args", args)
		return domain.User{}, err
	}

	result := user.ToDomain()
	u.logger.Debug(layer, method, "successfully completed", "updated_user", result)
	return result, nil
}

func (u *userRepo) DeleteUserByID(ctx context.Context, ID uuid.UUID) (domain.User, error) {
	const layer string = "repository"
	const method = "DeleteUserByID"

	u.logger.Debug(layer, method, "started", "id", ID)

	query, args, err := queries.BuildDeleteUserQuery(ID)
	if err != nil {
		u.logger.Error(layer, method, "failed to build query", err, "id", ID)
		return domain.User{}, err
	}

	u.logger.Debug(layer, method, "query built", "query", query, "args", args)

	var user dao.UserDAO
	err = u.q.QueryRow(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.Name,
		&user.Surname,
		&user.Nationality,
		&user.Patronymic,
		&user.Sex,
		&user.Age,
	)
	if err != nil {
		u.logger.Error(layer, method, "query execution failed", err, "query", query, "args", args)
		return domain.User{}, err
	}

	result := user.ToDomain()
	u.logger.Debug(layer, method, "successfully completed", "deleted_user", result)
	return result, nil
}
