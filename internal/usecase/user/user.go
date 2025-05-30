package user_usecase

import (
	"context"
	"database/sql"
	"errors"
	"sync"

	"github.com/FlyKarlik/effectiveMobile/internal/domain"
	user_drver "github.com/FlyKarlik/effectiveMobile/internal/driver/user"
	"github.com/FlyKarlik/effectiveMobile/internal/errs"
	user_repo "github.com/FlyKarlik/effectiveMobile/internal/repository/user"
	"github.com/FlyKarlik/effectiveMobile/pkg/generics"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
	"github.com/google/uuid"
	"golang.org/x/sync/errgroup"
)

type IUserUsecase interface {
	SearchUsers(ctx context.Context, pagination domain.Pagination, filter domain.UserFilter) generics.ItemsOutput[domain.User]
	CreateUser(ctx context.Context, input domain.CreateUserInput) error
	DeleteUserByID(ctx context.Context, ID uuid.UUID) error
	UpdateUserByID(ctx context.Context, ID uuid.UUID, input domain.UpdateUserInput) error
}

type userUsecase struct {
	logger     logger.Logger
	userRepo   user_repo.IUserRepository
	userDriver user_drver.IUserDriver
}

func New(logger logger.Logger, userRepo user_repo.IUserRepository, userDriver user_drver.IUserDriver) IUserUsecase {
	return &userUsecase{
		logger:     logger,
		userRepo:   userRepo,
		userDriver: userDriver,
	}
}

func (u *userUsecase) SearchUsers(ctx context.Context, pagination domain.Pagination, filter domain.UserFilter) generics.ItemsOutput[domain.User] {
	const layer = "usecase"
	const method = "SearchUsers"

	u.logger.Debug(layer, method, "started", "pagination", pagination, "filter", filter)

	var (
		count int64
		data  []domain.User
	)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		var err error
		count, err = u.userRepo.CountUsers(ctx, filter)
		if err != nil {
			u.logger.Error(layer, method, "failed to count users", err, "filter", filter)
			return err
		}
		u.logger.Debug(layer, method, "users counted", "count", count)
		return nil
	})

	g.Go(func() error {
		var err error
		data, err = u.userRepo.SearchUsers(ctx, pagination, filter)
		if err != nil {
			u.logger.Error(layer, method, "failed to search users", err, "pagination", pagination, "filter", filter)
			return err
		}
		u.logger.Debug(layer, method, "users fetched", "count", len(data))
		return nil
	})

	if err := g.Wait(); err != nil {
		u.logger.Error(layer, method, "operation failed", err)
		return generics.ItemsOutput[domain.User]{
			Success: false,
			Error:   err,
		}
	}

	u.logger.Debug(layer, method, "successfully completed", "total_count", count, "items_count", len(data))
	return generics.ItemsOutput[domain.User]{
		Success: true,
		Total:   count,
		Items:   data,
	}
}

func (u *userUsecase) DeleteUserByID(ctx context.Context, ID uuid.UUID) error {
	const layer = "usecase"
	const method = "DeleteUserByID"

	u.logger.Debug(layer, method, "started", "user_id", ID)

	deletedUser, err := u.userRepo.DeleteUserByID(ctx, ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.logger.Warn(layer, method, "user not found", err, "user_id", ID)
			return errs.ErrUserNotFound
		}
		u.logger.Error(layer, method, "failed to delete user", err, "user_id", ID)
		return errs.ErrUnknown
	}

	u.logger.Debug(layer, method, "successfully deleted user", "deleted_user", deletedUser)
	return nil
}

func (u *userUsecase) UpdateUserByID(ctx context.Context, ID uuid.UUID, input domain.UpdateUserInput) error {
	const layer = "usecase"
	const method = "UpdateUserByID"

	u.logger.Debug(layer, method, "started", "user_id", ID, "input", input)

	updatedUser, err := u.userRepo.UpdateUserByID(ctx, ID, input)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			u.logger.Warn(layer, method, "user not found", err, "user_id", ID)
			return errs.ErrUserNotFound
		}
		u.logger.Error(layer, method, "failed to update user", err, "user_id", ID, "input", input)
		return errs.ErrUnknown
	}

	u.logger.Debug(layer, method, "successfully updated user", "updated_user", updatedUser)
	return nil
}

func (u *userUsecase) CreateUser(ctx context.Context, input domain.CreateUserInput) error {
	const method = "CreateUser"
	const layer = "usecase"

	u.logger.Debug(layer, method, "started", "input", input)

	var (
		age         *int64
		nationality *string
		sex         *domain.SexEnum
	)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		u.logger.Debug(layer, method, "fetching age started", "name", input.Name)

		userAge, err := u.userDriver.GetUserAge(ctx, input.Name)
		if err != nil {
			u.logger.Error(layer, method, "Failed to get user age", err)
			return
		}

		if userAge != 0 {
			ageVal := int64(userAge)
			age = &ageVal
			u.logger.Debug(layer, method, "age fetched", "name", input.Name, "age", ageVal)
		} else {
			u.logger.Warn(layer, method, "age not found", nil, "name", input.Name)
		}
	}()

	go func() {
		defer wg.Done()
		u.logger.Debug(layer, method, "fetching nationality started", "name", input.Name)

		userNationality, err := u.userDriver.GetUserNationality(ctx, input.Name)
		if err != nil {
			u.logger.Error(layer, method, "Failed to get user nationality", err)
			return
		}

		if userNationality != "" {
			nationality = &userNationality
			u.logger.Debug(layer, method, "nationality fetched", "name", input.Name, "nationality", userNationality)
		} else {
			u.logger.Warn(layer, method, "nationality not found", nil, "name", input.Name)
		}
	}()

	go func() {
		defer wg.Done()
		u.logger.Debug(layer, method, "fetching sex started", "name", input.Name)

		userSex, err := u.userDriver.GetUserSex(ctx, input.Name)
		if err != nil {
			u.logger.Error(layer, method, "Failed to get user sex", err)
			return
		}

		if userSex != "" {
			sex = &userSex
			u.logger.Debug(layer, method, "sex fetched", "name", input.Name, "sex", userSex)
		} else {
			u.logger.Warn(layer, method, "sex not found", nil, "name", input.Name)
		}
	}()

	wg.Wait()

	u.logger.Debug(layer, method, "enrichment results",
		"age", age,
		"nationality", nationality,
		"sex", sex)

	input.Age = age
	input.Nationality = nationality
	input.Sex = sex

	u.logger.Debug(layer, method, "creating user in repository", "input", input)

	createdUser, err := u.userRepo.CreateUser(ctx, input)
	if err != nil {
		u.logger.Error(layer, method, "failed to create user", err, "input", input)
		return err
	}

	u.logger.Debug(layer, method, "user created successfully", "userID", createdUser)
	return nil
}
