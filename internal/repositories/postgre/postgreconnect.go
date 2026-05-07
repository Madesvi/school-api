// Package postgre provide connection to DB
package postgre

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable connect_timeout=3", host, user, password, dbname, port)
	slog.Debug("--- Try to connect to DB ---", "host", host, "port", port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		AllowGlobalUpdate: false, // Not allow delete or update all rows withot Where
		PrepareStmt:       true,  // Что бы не пересобирать SQL запросы при каждом вызове
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to DB: %w", err)
	}
	slog.Debug("--- Connected to DB ---", "host", host, "port", port)

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Failed to get database connection: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Ping to check DB
	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("Failed to connect to DB: %w", err)
	}

	return db, nil
}
