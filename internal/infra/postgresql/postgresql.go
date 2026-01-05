package postgresql

import (
	"fmt"
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(env *config.Env) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s",
		env.DBHost, env.DBUser, env.DBPassword, env.DBName, env.DBPort, env.DBSSLMode)

	logLevel := logger.Silent
	if env.AppEnv == "development" {
		logLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logLevel),
	})

	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		return nil, err
	}

	if env.AppEnv == "development" {
		db = db.Debug()
	}

	return db, nil
}
