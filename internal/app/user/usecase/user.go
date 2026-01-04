package usecase

import (
	"context"
	"log/slog"
	"math"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/infra/minio"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo  contract.UserRepositoryItf
	storyRepo contract.StoryRepositoryItf
	dictRepo  contract.DictionaryRepositoryItf
	storage   minio.MinioItf
	env       *config.Env
}

func NewUserUsecase(userRepo contract.UserRepositoryItf, storyRepo contract.StoryRepositoryItf, dictRepo contract.DictionaryRepositoryItf, storage minio.MinioItf, env *config.Env) contract.UserUsecaseItf {
	return &userUsecase{
		userRepo:  userRepo,
		storyRepo: storyRepo,
		dictRepo:  dictRepo,
		storage:   storage,
		env:       env,
	}
}

func (uc *userUsecase) GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, *response.APIError) {
	user, err := uc.userRepo.GetUserWithBadges(ctx, userID)
	if err != nil {
		slog.Error("failed to get user profile", "error", err)
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

	return &dto.UserProfileResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		CurrentTitle: string(user.CurrentTitle),
		Stats: dto.UserStatsResponse{
			TotalChapters:     totalChapters,
			CompletedChapters: user.LastChapterCompleted,
			ProgressPercent:   progressPercent,
			TotalVocabs:       totalVocabs,
			CollectedVocabs:   user.TotalWordsCollected,
		},
		Badges: badgeResponses,
	}, nil
}
