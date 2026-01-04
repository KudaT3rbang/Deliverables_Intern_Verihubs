package postgres

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/riverqueue/river/riverdriver/riverdatabasesql"
	"github.com/riverqueue/river/rivermigrate"
)

func RunMigration(pool *pgxpool.Pool, migrationsPath string) error {
	db := stdlib.OpenDBFromPool(pool)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db, migrationsPath); err != nil {
		return err
	}

	if err := runRiverMigrations(context.Background(), db); err != nil {
		return err
	}

	return nil
}

func runRiverMigrations(ctx context.Context, db *sql.DB) error {
	migrator, err := rivermigrate.New(riverdatabasesql.New(db), nil)
	if err != nil {
		return err
	}

	_, err = migrator.Migrate(ctx, rivermigrate.DirectionUp, nil)
	return err
}
