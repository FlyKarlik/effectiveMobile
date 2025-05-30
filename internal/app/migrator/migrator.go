package migrator

import (
	"database/sql"
	"errors"

	"github.com/FlyKarlik/effectiveMobile/config"
	"github.com/FlyKarlik/effectiveMobile/pkg/logger"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

var (
	ErrInvalidArgument = errors.New("usage up or down")
)

type AppMigrator struct {
	logger  logger.Logger
	cfg     *config.Config
	migrate *migrate.Migrate
}

func New(cfg *config.Config, logger logger.Logger) *AppMigrator {
	return &AppMigrator{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *AppMigrator) Migrate(migrateType string) error {
	const layer string = "migrator"
	const method string = "Migrate"

	a.logger.Info(layer, method, "Starting database migration", "type", migrateType)

	db, err := sql.Open("postgres", a.cfg.Infra.Postgres.ConnStr)
	if err != nil {
		a.logger.Error(layer, method, "Failed to open database connection", err)
		return err
	}
	defer func() {
		if err := db.Close(); err != nil {
			a.logger.Warn(layer, method, "Failed to close database connection", err)
		} else {
			a.logger.Debug(layer, method, "Database connection closed successfully")
		}
	}()

	a.logger.Debug(layer, method, "Creating postgres driver instance")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		a.logger.Error(layer, method, "Failed to create postgres driver", err)
		return err
	}

	migrationPath := "file://" + a.cfg.AppMigrator.MigrationsPath
	a.logger.Debug(layer, method, "Creating migrate instance",
		"migrationsPath", migrationPath)

	m, err := migrate.NewWithDatabaseInstance(
		migrationPath,
		"postgres",
		driver,
	)
	if err != nil {
		a.logger.Error(layer, method, "Failed to create migrate instance", err)
		return err
	}
	a.migrate = m

	a.logger.Info(layer, method, "Executing migration", "type", migrateType)
	switch migrateType {
	case "up":
		err = a.Up()
	case "down":
		err = a.down()
	default:
		a.logger.Error(layer, method, "Invalid migration type", ErrInvalidArgument,
			"providedType", migrateType)
		return ErrInvalidArgument
	}

	if err != nil {
		a.logger.Error(layer, method, "Migration failed", err)
		return err
	}

	a.logger.Info(layer, method, "Migration completed successfully")
	return nil
}

func (a *AppMigrator) down() error {
	const layer string = "migrator"
	const method string = "down"

	a.logger.Info(layer, method, "Executing DOWN migration")
	err := a.migrate.Down()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			a.logger.Info(layer, method, "No changes needed for DOWN migration")
			return nil
		}
		a.logger.Error(layer, method, "DOWN migration failed", err)
		return err
	}
	a.logger.Info(layer, method, "DOWN migration completed successfully")
	return nil
}

func (a *AppMigrator) Up() error {
	const layer string = "migrator"
	const method string = "Up"

	a.logger.Info(layer, method, "Executing UP migration")
	err := a.migrate.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			a.logger.Info(layer, method, "No changes needed for UP migration")
			return nil
		}
		a.logger.Error(layer, method, "UP migration failed", err)
		return err
	}
	a.logger.Info(layer, method, "UP migration completed successfully")
	return nil
}
