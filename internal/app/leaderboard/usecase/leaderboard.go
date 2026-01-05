package usecase

import (
	"context"
	"log/slog"

	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/infra/minio"
	"github.com/Ablebil/lathi-be/pkg/response"
)

type leaderboardUsecase struct {
	repo    contract.LeaderboardRepositoryItf
	storage minio.MinioItf
}

func NewLeaderboardUsecase(lbRepo contract.LeaderboardRepositoryItf, storage minio.MinioItf) contract.LeaderboardUsecaseItf {
	return &leaderboardUsecase{
		repo:    lbRepo,
		storage: storage,
	}
}

func (uc *leaderboardUsecase) GetLeaderboard(ctx context.Context) (*dto.LeaderboardResponse, *response.APIError) {
	entries, err := uc.repo.GetTopUsers(ctx, 5)
	if err != nil {
		slog.Error("failed to get top users", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	var topUsers []dto.LeaderboardItemResponse
	for _, entry := range entries {
		topUsers = append(topUsers, dto.LeaderboardItemResponse{
			Rank:     entry.Rank,
			UserID:   entry.UserID,
			Username: entry.Username,
			Title:    entry.Title,
			Score:    entry.Score,
		})
	}

	return &dto.LeaderboardResponse{
		TopUsers: topUsers,
	}, nil
}
