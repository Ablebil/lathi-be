package repository

import (
	"context"
	"fmt"

	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/internal/infra/redis"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type leaderboardRepository struct {
	db    *gorm.DB
	cache redis.RedisItf
}

func NewLeaderboardRepository(db *gorm.DB, cache redis.RedisItf) contract.LeaderboardRepositoryItf {
	return &leaderboardRepository{
		db:    db,
		cache: cache,
	}
}

func (r *leaderboardRepository) UpdateUserScore(ctx context.Context, userID uuid.UUID) error {
	var user entity.User
	err := r.db.WithContext(ctx).
		Select("last_chapter_completed", "total_words_collected").
		First(&user, userID).Error
	if err != nil {
		return err
	}

	score := calculateScore(user.LastChapterCompleted, user.TotalWordsCollected)
	return r.cache.ZAdd(ctx, "leaderboard:global", float64(score), userID.String())
}

func (r *leaderboardRepository) GetTopUsers(ctx context.Context, limit int) ([]contract.LeaderboardEntry, error) {
	results, err := r.cache.ZRevRangeWithScores(ctx, "leaderboard:global", 0, int64(limit-1))
	if err != nil {
		return nil, err
	}

	var entries []contract.LeaderboardEntry
	for i, result := range results {
		userID, err := uuid.Parse(fmt.Sprintf("%v", result.Member))
		if err != nil {
			continue
		}

		var user entity.User
		if err := r.db.WithContext(ctx).Select("username", "current_title").First(&user, userID).Error; err != nil {
			continue
		}

		entries = append(entries, contract.LeaderboardEntry{
			Rank:     i + 1,
			UserID:   userID,
			Username: user.Username,
			Title:    string(user.CurrentTitle),
			Score:    int(result.Score),
		})
	}

	return entries, nil
}

func (r *leaderboardRepository) GetUserRank(ctx context.Context, userID uuid.UUID) (int, int, error) {
	rank, err := r.cache.ZRevRank(ctx, "leaderboard:global", userID.String())
	if err != nil {
		return 0, 0, err
	}

	if rank == -1 {
		return 0, 0, nil // user not ranked yet
	}

	score, err := r.cache.ZScore(ctx, "leaderboard:global", userID.String())
	if err != nil {
		return 0, 0, err
	}

	return int(rank) + 1, int(score), nil
}

func (r *leaderboardRepository) RebuildLeaderboard(ctx context.Context) error {
	var users []entity.User
	err := r.db.WithContext(ctx).
		Where("is_verified = ?", true).
		Select("id", "last_chapter_completed", "total_words_collected").
		Find(&users).Error
	if err != nil {
		return err
	}

	for _, user := range users {
		score := calculateScore(user.LastChapterCompleted, user.TotalWordsCollected)
		if err := r.cache.ZAdd(ctx, "leaderboard:global", float64(score), user.ID.String()); err != nil {
			return err
		}
	}

	return nil
}

func calculateScore(chaptersCompleted, vocabsCollected int) int {
	return (chaptersCompleted * 100) + (vocabsCollected * 10)
}
