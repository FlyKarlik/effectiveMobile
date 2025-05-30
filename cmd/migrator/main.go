package main

import (
	"errors"
	"os"

	"github.com/FlyKarlik/effectiveMobile/config"
	"github.com/FlyKarlik/effectiveMobile/internal/app/migrator"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
	"github.com/FlyKarlik/effectiveMobile/pkg/validation"
)

var (
	ErrNotEnougthArgs = errors.New("not enouth arguments")
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		panic(ErrNotEnougthArgs)
	}

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := validation.Validate(cfg); err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.AppMigrator.LogLevel)
	if err != nil {
		panic(err)
	}

	migrator := migrator.New(cfg, logger)

	if err := migrator.Migrate(args[0]); err != nil {
		panic(err)
	}
}
