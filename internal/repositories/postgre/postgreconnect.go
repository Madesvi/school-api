// Package postgre provide connection to DB
package postgre

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBStorage struct {
	Gorm *gorm.DB
	Pgx  *pgxpool.Pool
}

func ConnectDB() (*DBStorage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable connect_timeout=3", host, user, password, dbname, port)
	slog.Debug("--- Try to connect to DB ---", "host", host, "port", port)
	gDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		AllowGlobalUpdate: false, // Not allow delete or update all rows withot Where
		PrepareStmt:       true,  // Что бы не пересобирать SQL запросы при каждом вызове
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}
	slog.Debug("--- Connected to DB ---", "host", host, "port", port)

	sqlDB, err := gDB.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping to check DB
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	config.MaxConns = 50
	config.MinConns = 5
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = 15 * time.Minute

	pPool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgxpool: %w", err)
	}
	err = pPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping pgx: %w", err)
	}

	return &DBStorage{Gorm: gDB, Pgx: pPool}, nil
}
