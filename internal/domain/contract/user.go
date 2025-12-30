package contract

import (
	"context"

	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepositoryItf interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	CreateUser(ctx context.Context, user *entity.User) error
	UpdateUser(ctx context.Context, user *entity.User) error
}
