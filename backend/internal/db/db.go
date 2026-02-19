package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const defaultDSN = "postgres://exchange:exchange_local@localhost:5432/exchange_dev?sslmode=disable"

// Connect opens a GORM connection to PostgreSQL. If dsn is empty, uses the local Docker default.
func Connect(dsn string) (*gorm.DB, error) {
	if dsn == "" {
		dsn = defaultDSN
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		return nil, fmt.Errorf("gorm.Open: %w", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("db.DB(): %w", err)
	}
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping: %w", err)
	}
	return db, nil
}
