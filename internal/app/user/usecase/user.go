package usecase

import (
	"context"
	"fmt"
	"log/slog"
	"math"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/infra/minio"
	"github.com/Ablebil/lathi-be/internal/infra/redis"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo  contract.UserRepositoryItf
	storyRepo contract.StoryRepositoryItf
	dictRepo  contract.DictionaryRepositoryItf
	lbRepo    contract.LeaderboardRepositoryItf
	storage   minio.MinioItf
	cache     redis.RedisItf
	env       *config.Env
}

func NewUserUsecase(userRepo contract.UserRepositoryItf, storyRepo contract.StoryRepositoryItf, dictRepo contract.DictionaryRepositoryItf, lbRepo contract.LeaderboardRepositoryItf, storage minio.MinioItf, cache redis.RedisItf, env *config.Env) contract.UserUsecaseItf {
	return &userUsecase{
		userRepo:  userRepo,
		storyRepo: storyRepo,
		dictRepo:  dictRepo,
		lbRepo:    lbRepo,
		storage:   storage,
		cache:     cache,
		env:       env,
	}
}

func (uc *userUsecase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, *response.APIError) {
	user, err := uc.userRepo.GetUserWithBadges(ctx, userID)
	if err != nil {
		return nil, response.ErrNotFound("Akun ga ditemukan, coba daftar dulu ya")
	}

	totalChapters, err := uc.storyRepo.CountChapters(ctx)
	if err != nil {
		slog.Error("failed to count chapters", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	totalVocabs, err := uc.dictRepo.CountTotalVocabs(ctx)
	if err != nil {
		slog.Error("failed to count vocabs", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	progressPercent := 0.0
	if totalChapters > 0 {
		progressPercent = (float64(user.LastChapterCompleted) / float64(totalChapters)) * 100
		if progressPercent > 100 {
			progressPercent = 100
		}
	}
	progressPercent = math.Round(progressPercent*100) / 100

	var badgeResponses []dto.UserBadgeResponse
	for _, ub := range user.UserBadges {
		badgeResponses = append(badgeResponses, dto.UserBadgeResponse{
			Name:        ub.Badge.Name,
			Description: ub.Badge.Description,
			IconURL:     uc.storage.GetObjectURL(ub.Badge.IconURL),
			EarnedAt:    ub.EarnedAt,
		})
	}

	var lbInfo *dto.UserLeaderboardInfoResponse
	rank, score, err := uc.lbRepo.GetUserRank(ctx, userID)
	if err != nil {
		slog.Error("failed to get user rank", "error", err)
	} else if rank > 0 {
		lbInfo = &dto.UserLeaderboardInfoResponse{
			Rank:  rank,
			Score: score,
		}
	}

	return &dto.UserProfileResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		AvatarURL:    uc.storage.GetObjectURL(user.AvatarURL),
		CurrentTitle: string(user.CurrentTitle),
		Stats: dto.UserStatsResponse{
			TotalChapters:     totalChapters,
			CompletedChapters: user.LastChapterCompleted,
			ProgressPercent:   progressPercent,
			TotalVocabs:       totalVocabs,
			CollectedVocabs:   user.TotalWordsCollected,
		},
		Badges:          badgeResponses,
		LeaderboardInfo: lbInfo,
	}, nil
}

func (uc *userUsecase) EditUserProfile(ctx context.Context, userID uuid.UUID, req *dto.EditUserProfileRequest) (*dto.UserProfileResponse, *response.APIError) {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	if user.Username == req.Username {
		return nil, response.ErrBadRequest("Username baru sama dengan yang lama")
	}

	existingUser, err := uc.userRepo.GetUserByUsername(ctx, req.Username)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}
	if existingUser != nil {
		return nil, response.ErrConflict("Username ini udah dipake, coba yang lain ya")
	}

	user.Username = req.Username
	if err := uc.userRepo.UpdateUser(ctx, user); err != nil {
		slog.Error("failed to update user", "error", err)
		return nil, response.ErrInternal("Coba lagi nanti ya!")
	}

	return uc.GetUserProfile(ctx, userID)
}

func (uc *userUsecase) DeleteAccount(ctx context.Context, userID uuid.UUID, refreshToken string) *response.APIError {
	user, err := uc.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		slog.Error("failed to get user", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}
	if user == nil {
		return response.ErrNotFound("Akun ga ditemukan, coba daftar dulu ya")
	}

	if err := uc.userRepo.DeleteUser(ctx, userID); err != nil {
		slog.Error("failed to delete user", "error", err)
		return response.ErrInternal("Coba lagi nanti ya!")
	}

	if refreshToken != "" {
		cacheKey := fmt.Sprintf("refresh:%s", refreshToken)
		if err := uc.cache.Del(ctx, cacheKey); err != nil {
			slog.Warn("failed to delete refresh token", "error", err)
		}
	}

	if err := uc.lbRepo.RemoveUserFromLeaderboard(ctx, userID); err != nil {
		slog.Warn("failed to remove user from leaderboard", "error", err)
	}

	return nil
}
