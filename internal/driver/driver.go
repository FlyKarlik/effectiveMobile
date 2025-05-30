package driver

import (
	user_drver "github.com/FlyKarlik/effectiveMobile/internal/driver/user"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
)

type Driver struct {
	user_drver.IUserDriver
}

type driverOptions func(r *Driver) error

func New(cfgs ...driverOptions) (*Driver, error) {
	os := &Driver{}

	for _, cfg := range cfgs {
		err := cfg(os)
		if err != nil {
			return nil, err
		}
	}

	return os, nil
}

func WithUserDriver(logger logger.Logger) driverOptions {
	return func(r *Driver) error {
		r.IUserDriver = user_drver.New(logger)
		return nil
	}
}
