package seed

import (
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/pkg/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSeeder struct{}

func (s *UserSeeder) Run(db *gorm.DB) error {
	slog.Info("seeding user domain...")

	hasher := bcrypt.NewBcrypt()
	password, err := hasher.Hash("Str0ngP@ssw0rD")
	if err != nil {
		return err
	}

	users := []entity.User{
		{
			ID:           uuid.New(),
			Username:     "valen",
			Email:        "valen@lathi.id",
			Password:     password,
			IsVerified:   true,
			CurrentTitle: entity.Priyayi,
		},
		{
			ID:           uuid.New(),
			Username:     "soma",
			Email:        "soma@lathi.id",
			Password:     password,
			IsVerified:   true,
			CurrentTitle: entity.Cantrik,
		},
	}

	for _, u := range users {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"username", "password", "is_verified", "current_title"}),
		}).Create(&u).Error; err != nil {
			slog.Error("failed to seed user", "email", u.Email, "error", err)
			return err
		}
	}

	slog.Info("user seeding completed successfully")
	return nil
}
