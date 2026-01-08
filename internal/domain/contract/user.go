package contract

import (
	"context"
	"time"

	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type UserUsecaseItf interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, *response.APIError)
	EditUserProfile(ctx context.Context, userID uuid.UUID, req *dto.EditUserProfileRequest) (*dto.UserProfileResponse, *response.APIError)
	DeleteAccount(ctx context.Context, userID uuid.UUID, refreshToken string) *response.APIError
}

type UserRepositoryItf interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserWithBadges(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	GetUserLastCompletedChapter(ctx context.Context, userID uuid.UUID) (int, error)
	UpdateUserLastCompletedChapter(ctx context.Context, userID uuid.UUID, orderIndex int) error
	IncrementUserWordCount(ctx context.Context, userID uuid.UUID, amount int) error
	UpdateUserTitle(ctx context.Context, userID uuid.UUID, title entity.Title) error
	AssignBadge(ctx context.Context, userID uuid.UUID, badgeCode string) error
	DeleteUnverifiedUsers(ctx context.Context, threshold time.Time) (int64, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}
