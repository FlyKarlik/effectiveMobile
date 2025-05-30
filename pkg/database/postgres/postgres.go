package postgres

import (
	"context"
	"fmt"

	"github.com/FlyKarlik/effectiveMobile/config"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Querier interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	Begin(ctx context.Context) (pgx.Tx, error)
	Ping(ctx context.Context) error
	Close()
}

func NewPostgresDB(config *config.PostgreSQL) (Querier, error) {
	dataSourceName := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		config.Host,
		config.Port,
		config.User,
		config.Database,
		config.Password,
	)

	pool, err := pgxpool.New(context.Background(), dataSourceName)
	if err != nil {
		return nil, err
	}

	pool.Config().AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	if err = pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
