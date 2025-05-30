package main

import (
	"github.com/FlyKarlik/effectiveMobile/config"
	"github.com/FlyKarlik/effectiveMobile/internal/app/users"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
	validator "github.com/FlyKarlik/effectiveMobile/pkg/validation"
)

// @title Users API
// @version 1.0
// @description API documentation for the Users backend service.

// @contact.name API Support
// @contact.url https://github.com/FlyKarlik
// @contact.email nikitasavin191@gmail.com

// @host localhost:8000
// @BasePath /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	if err := validator.Validate(cfg); err != nil {
		panic(err)
	}

	logger, err := logger.New(cfg.AppUsers.LogLevel)
	if err != nil {
		panic(err)
	}

	usersApp := users.New(logger, cfg)
	if err := usersApp.Start(); err != nil {
		panic(err)
	}
}
