package handler

import (
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/middleware"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userHandler struct {
	uc contract.UserUsecaseItf
}

func NewUserHandler(router fiber.Router, mw middleware.MiddlewareItf, uc contract.UserUsecaseItf) {
	handler := &userHandler{
		uc: uc,
	}

	userRouter := router.Group("/users", mw.Authenticate)
	userRouter.Get("/profile", handler.getProfile)
}

func (h *userHandler) getProfile(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	resp, apiErr := h.uc.GetUserProfile(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Profilmu berhasil dimuat", resp)
}
