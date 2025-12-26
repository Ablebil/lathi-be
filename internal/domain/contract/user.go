package contract

import "github.com/Ablebil/lathi-be/internal/domain/entity"

type UserRepositoryItf interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByRefreshToken(token string) (*entity.User, error)
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
}
