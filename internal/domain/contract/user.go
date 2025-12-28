package contract

import (
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"github.com/google/uuid"
)

type UserRepositoryItf interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(id uuid.UUID) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
}
