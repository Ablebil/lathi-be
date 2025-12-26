package repository

import (
	"errors"

	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/entity"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) contract.UserRepositoryItf {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) GetUserByRefreshToken(token string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("refresh_token = ?", token).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *userRepository) CreateUser(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) UpdateUser(user *entity.User) error {
	return r.db.Save(user).Error
}
