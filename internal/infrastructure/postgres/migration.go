package postgres

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func RunMigration(pool *pgxpool.Pool, migrationsDir string) error {
	conn := stdlib.OpenDBFromPool(pool)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(conn, migrationsDir); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	version, err := goose.GetDBVersion(conn)
	if err != nil {
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	log.Printf("Successfully applied migrations. Current version: %d", version)
	return nil
}

func RollbackMigration(pool *pgxpool.Pool, migrationsDir string) error {
	connConfig := pool.Config().ConnConfig
	connString := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		connConfig.Host,
		connConfig.Port,
		connConfig.User,
		connConfig.Password,
		connConfig.Database,
	)

	db, err := sql.Open("pgx", connString)
	if err != nil {
		return fmt.Errorf("failed to open database for rollback: %w", err)
	}

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Down(db, migrationsDir); err != nil {
		return fmt.Errorf("failed to rollback migration: %w", err)
	}

	log.Println("Successfully rolled back migration")
	return nil
}
