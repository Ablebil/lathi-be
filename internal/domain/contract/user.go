package contract

import (
	"context"

	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/google/uuid"
)

type UserUsecaseItf interface {
	GetUserProfile(ctx context.Context, userID uuid.UUID) (*dto.UserProfileResponse, *response.APIError)
}

type UserRepositoryItf interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	GetUserWithBadges(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
	GetUserLastCompletedChapter(ctx context.Context, userID uuid.UUID) (int, error)
	UpdateUserLastCompletedChapter(ctx context.Context, userID uuid.UUID, orderIndex int) error
	IncrementUserWordCount(ctx context.Context, userID uuid.UUID, amount int) error
	UpdateUserTitle(ctx context.Context, userID uuid.UUID, title entity.Title) error
}
