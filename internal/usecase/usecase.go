package usecase

import (
	user_drver "github.com/FlyKarlik/effectiveMobile/internal/driver/user"
	user_repo "github.com/FlyKarlik/effectiveMobile/internal/repository/user"
	user_usecase "github.com/FlyKarlik/effectiveMobile/internal/usecase/user"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
)

type Usecase struct {
	user_usecase.IUserUsecase
}

type repoOptions func(r *Usecase) error

func New(cfgs ...repoOptions) (*Usecase, error) {
	os := &Usecase{}

	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}

	return os, nil
}

func WithUserUsecase(logger logger.Logger, userRepo user_repo.IUserRepository, userDriver user_drver.IUserDriver) repoOptions {
	return func(r *Usecase) error {
		r.IUserUsecase = user_usecase.New(logger, userRepo, userDriver)
		return nil
	}
}
