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
		return response.Error(ctx, response.ErrBadRequest("Data yang kamu kirim belum pas, coba cek lagi ya"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if apiErr := h.uc.Register(ctx.Context(), req); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusCreated, "Pendaftaran berhasil, cek email kamu buat verifikasi, ya", nil)
}

func (h *authHandler) verify(ctx *fiber.Ctx) error {
	req := new(dto.VerifyRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Data yang kamu kirim belum pas, coba cek lagi ya"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if apiErr := h.uc.Verify(ctx.Context(), req); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Email kamu udah diverifikasi, yuk login sekarang!", nil)
}

func (h *authHandler) login(ctx *fiber.Ctx) error {
	req := new(dto.LoginRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Data yang kamu kirim belum pas, coba cek lagi ya"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	resp, apiErr := h.uc.Login(ctx.Context(), req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Login sukses! Yuk mulai eksplorasi!", resp)
}

func (h *authHandler) refresh(ctx *fiber.Ctx) error {
	req := new(dto.RefreshRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Data yang kamu kirim belum pas, coba cek lagi ya"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	resp, apiErr := h.uc.Refresh(ctx.Context(), req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Sesi kamu udah diperbarui, yuk lanjut eksplorasi!", resp)
}

func (h *authHandler) logout(ctx *fiber.Ctx) error {
	req := new(dto.LogoutRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Data yang kamu kirim belum pas, coba cek lagi ya"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	if apiErr := h.uc.Logout(ctx.Context(), req); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Logout berhasil, sampai jumpa lagi!", nil)
}
