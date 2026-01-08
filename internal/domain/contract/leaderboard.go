package contract

import (
	"context"

	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type LeaderboardUsecaseItf interface {
	GetLeaderboard(ctx context.Context) (*dto.LeaderboardResponse, *response.APIError)
}

type LeaderboardRepositoryItf interface {
	UpdateUserScore(ctx context.Context, userID uuid.UUID) error
	GetTopUsers(ctx context.Context, limit int) ([]LeaderboardEntry, error)
	GetUserRank(ctx context.Context, userID uuid.UUID) (rank int, score int, err error)
	RebuildLeaderboard(ctx context.Context) error
	RemoveUserFromLeaderboard(ctx context.Context, userID uuid.UUID) error
}

type LeaderboardEntry struct {
	Rank      int
	UserID    uuid.UUID
	Username  string
	AvatarURL string
	Title     string
	Score     int
}
