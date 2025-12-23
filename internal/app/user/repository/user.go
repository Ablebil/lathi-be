package repository

import "gorm.io/gorm"

type UserRepositoryItf interface{}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepositoryItf {
	return &UserRepository{db}
}
