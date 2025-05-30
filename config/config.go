package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/joho/godotenv/autoload"
)

type Config struct {
	AppUsers    AppUsers       `validate:"required"`
	AppMigrator AppMigrator    `validate:"required"`
	Infra       Infrastructure `validate:"required"`
}

type AppUsers struct {
	LogLevel string `env:"APP__USERS__LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	AppPort  string `env:"APP__USERS__PORT" validate:"required,numeric,min=4,max=5"`
	AppHost  string `env:"APP__USERS__HOST" validate:"required,hostname_rfc1123|ipv4|ipv6"`
}

type AppMigrator struct {
	LogLevel       string `env:"APP__MIGRATOR__LOG_LEVEL" validate:"required,oneof=debug info warn error"`
	MigrationsPath string `env:"APP__MIGRATOR__MIGRATIONS_PATH" validate:"required"`
}

type Infrastructure struct {
	Postgres PostgreSQL `validate:"required"`
}

type PostgreSQL struct {
	Host     string `env:"INFRA__POSTGRES__HOST" validate:"required,hostname|ip"`
	Port     string `env:"INFRA__POSTGRES__PORT" validate:"required,numeric"`
	User     string `env:"INFRA__POSTGRES__USER" validate:"required"`
	Password string `env:"INFRA__POSTGRES__PASSWORD" validate:"required"`
	Database string `env:"INFRA__POSTGRES__DATABASE" validate:"required,alphaunicode"`
	ConnStr  string `env:"INFRA__POSTGRES__CONN_STR" validate:"required,url"`
}

func New() (*Config, error) {
	cfg := &Config{}
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
