package usecase

import "github.com/Ablebil/lathi-be/internal/app/user/repository"

type UserUsecaseItf interface{}

type UserUsecase struct {
	userRepository repository.UserRepositoryItf
}

func NewUserUsecase(r repository.UserRepositoryItf) UserUsecaseItf {
	return &UserUsecase{
		userRepository: r,
	}
}
