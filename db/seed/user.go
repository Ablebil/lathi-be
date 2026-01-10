package seed

import (
	"context"
	"log/slog"
	"time"

	lbRepo "github.com/Ablebil/lathi-be/internal/app/leaderboard/repository"
	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/infra/redis"
	"github.com/Ablebil/lathi-be/pkg/bcrypt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserSeeder struct {
	env *config.Env
}

func NewUserSeeder(env *config.Env) *UserSeeder {
	return &UserSeeder{env: env}
}

func (s *UserSeeder) Run(db *gorm.DB) error {
	slog.Info("seeding user domain...")

	hasher := bcrypt.NewBcrypt()
	password, err := hasher.Hash("Str0ngP@ssw0rD")
	if err != nil {
		return err
	}

	avatarURL := "avatars/default.webp"

	var badges []entity.Badge
	if err := db.Find(&badges).Error; err != nil {
		return err
	}
	badgeMap := make(map[string]uuid.UUID)
	for _, b := range badges {
		badgeMap[b.Code] = b.ID
	}

	users := []entity.User{
		{
			ID:                   uuid.New(),
			Username:             "valen",
			Email:                "valen@lathi.id",
			Password:             password,
			AvatarURL:            avatarURL,
			IsVerified:           true,
			CurrentTitle:         entity.Priyayi,
			LastChapterCompleted: 4,
			TotalWordsCollected:  58,
		},
		{
			ID:                   uuid.New(),
			Username:             "soma",
			Email:                "soma@lathi.id",
			Password:             password,
			AvatarURL:            avatarURL,
			IsVerified:           true,
			CurrentTitle:         entity.Abdi,
			LastChapterCompleted: 2,
			TotalWordsCollected:  25,
		},
		{
			ID:                   uuid.New(),
			Username:             "laras",
			Email:                "laras@lathi.id",
			Password:             password,
			AvatarURL:            avatarURL,
			IsVerified:           true,
			CurrentTitle:         entity.Cantrik,
			LastChapterCompleted: 1,
			TotalWordsCollected:  15,
		},
		{
			ID:                   uuid.New(),
			Username:             "budi",
			Email:                "budi@lathi.id",
			Password:             password,
			AvatarURL:            avatarURL,
			IsVerified:           true,
			CurrentTitle:         entity.Cantrik,
			LastChapterCompleted: 0,
			TotalWordsCollected:  7,
		},
		{
			ID:                   uuid.New(),
			Username:             "sari",
			Email:                "sari@lathi.id",
			Password:             password,
			AvatarURL:            avatarURL,
			IsVerified:           true,
			CurrentTitle:         entity.Cantrik,
			LastChapterCompleted: 1,
			TotalWordsCollected:  12,
		},
	}

	for _, u := range users {
		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "email"}},
			DoUpdates: clause.AssignmentColumns([]string{"username", "password", "is_verified", "current_title", "last_chapter_completed", "total_words_collected"}),
		}).Create(&u).Error; err != nil {
			slog.Error("failed to seed user", "email", u.Email, "error", err)
			return err
		}

		var userBadges []entity.UserBadge
		addBadge := func(code string) {
			if id, ok := badgeMap[code]; ok {
				userBadges = append(userBadges, entity.UserBadge{
					UserID:   u.ID,
					BadgeID:  id,
					EarnedAt: time.Now(),
				})
			}
		}

		switch u.Username {
		case "valen":
			addBadge("ch1_completion")
			addBadge("vocab_collector_1")
			addBadge("perfect_heart")
			addBadge("all_chapters_completion")
		case "soma":
			addBadge("ch1_completion")
			addBadge("vocab_collector_1")
		case "laras":
			addBadge("ch1_completion")
		case "sari":
			addBadge("ch1_completion")
			addBadge("perfect_heart")
		}

		if len(userBadges) > 0 {
			if err := db.Clauses(clause.OnConflict{
				DoNothing: true,
			}).Create(&userBadges).Error; err != nil {
				slog.Error("failed to seed user badges", "username", u.Username, "error", err)
				return err
			}
		}
	}

	slog.Info("user seeding completed successfully")

	// populate leaderboard in redis
	slog.Info("populating leaderboard in redis")
	cache := redis.New(s.env)
	lbRepository := lbRepo.NewLeaderboardRepository(db, cache)

	ctx := context.Background()
	if err := lbRepository.RebuildLeaderboard(ctx); err != nil {
		slog.Error("failed to populate leaderboard", "error", err)
		return err
	}
	slog.Info("leaderboard populated successfully")

	return nil
}
