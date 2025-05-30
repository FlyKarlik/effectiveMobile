package repository

import (
	user_repo "github.com/FlyKarlik/effectiveMobile/internal/repository/user"
	"github.com/FlyKarlik/effectiveMobile/pkg/database/postgres"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
)

type Repository struct {
	user_repo.IUserRepository
}

type repoOptions func(r *Repository) error

func New(cfgs ...repoOptions) (*Repository, error) {
	os := &Repository{}

	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}

	return os, nil
}

func WithUserRepo(logger logger.Logger, q postgres.Querier) repoOptions {
	return func(r *Repository) error {
		r.IUserRepository = user_repo.New(logger, q)
		return nil
	}
}
