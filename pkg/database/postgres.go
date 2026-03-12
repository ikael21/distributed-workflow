package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type PostgresConfig struct {
	DSN             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type postgres struct {
	DB *bun.DB
}

func NewPostgres(cfg PostgresConfig) (*postgres, error) {
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(cfg.DSN)))
	sqldb.SetMaxOpenConns(cfg.MaxOpenConns)
	sqldb.SetMaxIdleConns(cfg.MaxIdleConns)
	sqldb.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	db := bun.NewDB(sqldb, pgdialect.New())
	err := db.PingContext(context.Background())

	return &postgres{DB: db}, err
}

func (p *postgres) Close(_ context.Context) error {
	return p.DB.Close()
}
