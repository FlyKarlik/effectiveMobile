package users

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"

	"github.com/FlyKarlik/effectiveMobile/config"
	http_handler "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/handler"
	http_middleware "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/middleware"
	http_router "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/router"
	http_server "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/server"
	"github.com/FlyKarlik/effectiveMobile/internal/driver"
	"github.com/FlyKarlik/effectiveMobile/internal/repository"
	"github.com/FlyKarlik/effectiveMobile/internal/usecase"
	"github.com/FlyKarlik/effectiveMobile/pkg/database/postgres"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"
)

type AppUsers struct {
	logger logger.Logger
	cfg    *config.Config
}

func New(logger logger.Logger, cfg *config.Config) *AppUsers {
	return &AppUsers{
		logger: logger,
		cfg:    cfg,
	}
}

func (a *AppUsers) Start() error {
	const layer = "users"
	const method = "Start"

	a.logger.Info(layer, method, "Starting users application")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a.logger.Info(layer, method, "Connecting to database")
	dbConn, err := postgres.NewPostgresDB(&a.cfg.Infra.Postgres)
	if err != nil {
		a.logger.Error(layer, method, "Failed to connect to database", err)
		return err
	}
	defer func() {
		a.logger.Debug(layer, method, "Closing database connection")
		dbConn.Close()
		a.logger.Info(layer, method, "Database connection closed")
	}()

	a.logger.Info(layer, method, "Initializing repository")
	repo, err := repository.New(
		repository.WithUserRepo(a.logger, dbConn),
	)
	if err != nil {
		a.logger.Error(layer, method, "Failed to initialize repository", err)
		return err
	}

	driver, err := driver.New(
		driver.WithUserDriver(a.logger),
	)
	if err != nil {
		a.logger.Error(layer, method, "Failed to initialize driver", err)
		return err
	}

	a.logger.Info(layer, method, "Initializing usecase")
	usecase, err := usecase.New(
		usecase.WithUserUsecase(a.logger, repo.IUserRepository, driver.IUserDriver),
	)
	if err != nil {
		a.logger.Error(layer, method, "Failed to initialize usecase", err)
		return err
	}

	a.logger.Info(layer, method, "Initializing HTTP components")
	httpHandler := http_handler.New(a.logger, usecase)
	httpMiddleware := http_middleware.New()
	httpRouter := http_router.New(httpMiddleware, httpHandler)
	httpServer := http_server.New(a.cfg, httpRouter)

	a.logger.Info(layer, method, "Starting HTTP server",
		"port", a.cfg.AppUsers.AppPort,
		"host", a.cfg.AppUsers.AppHost)

	go func() {
		if err := httpServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error(layer, method, "HTTP server fatal error", err)
			os.Exit(1)
		}
	}()

	a.logger.Debug(layer, method, "Setting up graceful shutdown handler")
	a.signalHandler(ctx)

	a.logger.Info(layer, method, "Shutting down HTTP server")
	if err := httpServer.Shuttdown(ctx); err != nil {
		a.logger.Error(layer, method, "Failed to shutdown HTTP server", err)
		return err
	}
	a.logger.Info(layer, method, "HTTP server shutdown completed")

	return nil
}

func (a *AppUsers) signalHandler(ctx context.Context) {
	const layer = "users"
	const method = "signalHandler"

	a.logger.Debug(layer, method, "Setting up signal handler")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	defer func() {
		signal.Stop(signalChan)
		close(signalChan)
		a.logger.Debug(layer, method, "Signal handler cleaned up")
	}()

	select {
	case <-ctx.Done():
		a.logger.Debug(layer, method, "Context cancelled, exiting signal handler")
		return
	case sig := <-signalChan:
		a.logger.Info(layer, method, "Received signal, shutting down", "signal", sig)
		return
	}
}
