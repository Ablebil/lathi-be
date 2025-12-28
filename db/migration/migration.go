package migration

import (
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/infra/postgresql"
)

func Migrate(env *config.Env, action string) {
	db, err := postgresql.New(env)
	if err != nil {
		slog.Error("failed to connect to database for migration", "error", err)
		return
	}

	models := []interface{}{
		&entity.User{},
		&entity.Dictionary{},
		&entity.UserVocabulary{},
		&entity.Chapter{},
		&entity.Slide{},
		&entity.UserStorySession{},
	}

	switch action {
	case "up":
		if err := db.AutoMigrate(models...); err != nil {
			slog.Error("migration failed", "error", err)
		}
	case "down":
		if err := db.Migrator().DropTable(models...); err != nil {
			slog.Error("migration rollback failed", "error", err)
		}
	}

	slog.Info("migration done")
}
