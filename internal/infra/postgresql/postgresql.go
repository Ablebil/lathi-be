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

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Default.LogMode(logger.Info),
	})

	db.Debug()

	if err != nil {
		slog.Error("failed to connect to postgres", "error", err)
		return nil, err
	}

	return db, nil
}
