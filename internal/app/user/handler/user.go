package handler

import (
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/middleware"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/Ablebil/lathi-be/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type userHandler struct {
	val validator.ValidatorItf
	uc  contract.UserUsecaseItf
}

func NewUserHandler(router fiber.Router, validator validator.ValidatorItf, mw middleware.MiddlewareItf, userUc contract.UserUsecaseItf) {
	handler := &userHandler{
		val: validator,
		uc:  userUc,
	}

	userRouter := router.Group("/users", mw.Authenticate)
	userRouter.Get("/profile", handler.getProfile)
	userRouter.Patch("/profile", handler.editProfile)
	userRouter.Delete("/account", handler.deleteAccount)
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

func (h *userHandler) editProfile(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	req := new(dto.EditUserProfileRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Data yang kamu kirim belum pas, coba cek lagi ya"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	resp, apiErr := h.uc.EditUserProfile(ctx.Context(), userID, req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Profilmu berhasil diperbarui", resp)
}

func (h *userHandler) deleteAccount(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	refreshToken := ctx.Cookies("refresh_token")

	if apiErr := h.uc.DeleteAccount(ctx.Context(), userID, refreshToken); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	ctx.Cookie(&fiber.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HTTPOnly: true,
		Secure:   true,
		SameSite: fiber.CookieSameSiteNoneMode,
	})

	return ctx.SendStatus(fiber.StatusNoContent)
}
