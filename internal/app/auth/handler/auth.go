package handler

import (
	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/Ablebil/lathi-be/pkg/validator"
	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	val validator.ValidatorItf
	uc  contract.AuthUsecaseItf
}

func NewAuthHandler(router fiber.Router, validator validator.ValidatorItf, authUc contract.AuthUsecaseItf) {
	handler := authHandler{
		val: validator,
		uc:  authUc,
	}

	authRouter := router.Group("/auth")
	authRouter.Post("/register", handler.register)
	authRouter.Post("/verify", handler.verify)
	authRouter.Post("/login", handler.login)
	authRouter.Post("/refresh", handler.refresh)
	authRouter.Post("/logout", handler.logout)
}

func (h *authHandler) register(ctx *fiber.Ctx) error {
	req := new(dto.RegisterRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("failed parsing request body"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if err := h.uc.Register(ctx.Context(), req); err != nil {
		return response.Error(ctx, err, nil)
	}

	return response.Success(ctx, fiber.StatusCreated, "registration successful, please check your email for verification", nil)
}

func (h *authHandler) verify(ctx *fiber.Ctx) error {
	req := new(dto.VerifyRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("failed parsing request body"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if err := h.uc.Verify(ctx.Context(), req); err != nil {
		return response.Error(ctx, err, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "email verification successful", nil)
}

func (h *authHandler) login(ctx *fiber.Ctx) error {
	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("failed to parsing request body"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	accessToken, refreshToken, err := h.uc.Login(ctx.Context(), req)
	if err != nil {
		return response.Error(ctx, err, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "login successful", dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *authHandler) refresh(ctx *fiber.Ctx) error {
	req := new(dto.RefreshRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("failed to parsing request body"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	accessToken, refreshToken, err := h.uc.Refresh(ctx.Context(), req)
	if err != nil {
		return response.Error(ctx, err, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "refresh tokens successful", dto.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *authHandler) logout(ctx *fiber.Ctx) error {
	req := new(dto.LogoutRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("failed to parsing request body"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if err := h.uc.Logout(ctx.Context(), req); err != nil {
		return response.Error(ctx, err, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "logout successful", nil)
}
