package handler

import (
	"github.com/Ablebil/lathi-be/internal/app/user/usecase"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	userUsecase usecase.UserUsecaseItf
}

func NewUserHandler(r fiber.Router, uc usecase.UserUsecaseItf) {
}
