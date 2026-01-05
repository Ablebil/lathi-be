package seed

import (
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/infra/postgresql"
	"gorm.io/gorm"
)

type Seeder interface {
	Run(db *gorm.DB) error
}

var registry = map[string]Seeder{
	"badge": &BadgeSeeder{},
	"story": &StorySeeder{},
}

var executionOrder = []string{"badge", "user", "story"}

func Seed(env *config.Env, domain string) {
	db, err := postgresql.New(env)
	if err != nil {
		slog.Error("failed to connect to database for seeding", "error", err)
		return
	}

	slog.Info("starting database seeding...")

	// specific domain seeding
	if domain != "" {
		var seeder Seeder
		if domain == "user" {
			seeder = NewUserSeeder(env)
		} else {
			var ok bool
			seeder, ok = registry[domain]
			if !ok {
				slog.Error("seeder not found for domain", "domain", domain)
				return
			}
		}

		slog.Info("running specific seeder", "domain", domain)
		if err := seeder.Run(db); err != nil {
			slog.Error("failed to execute seeder", "domain", domain, "error", err)
		} else {
			slog.Info("seeder completed", "domain", domain)
		}
		return
	}

	// run all seeders
	slog.Info("running all seeders...")
	for _, name := range executionOrder {
		var seeder Seeder
		if name == "user" {
			seeder = NewUserSeeder(env)
		} else {
			seeder = registry[name]
		}

		slog.Info("running seeder", "domain", name)

		if err := seeder.Run(db); err != nil {
			slog.Error("failed to run seeder", "domain", name, "error", err)
		}
	}

	slog.Info("all seeders execution process finished")
}
