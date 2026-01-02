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

type storyHandler struct {
	val validator.ValidatorItf
	uc  contract.StoryUsecaseItf
}

func NewStoryHandler(router fiber.Router, validator validator.ValidatorItf, mw middleware.MiddlewareItf, storyUc contract.StoryUsecaseItf) {
	handler := storyHandler{
		val: validator,
		uc:  storyUc,
	}

	storyRouter := router.Group("/stories", mw.Authenticate)
	storyRouter.Get("/chapters", handler.getChapterList)
	storyRouter.Get("/chapters/:id/content", handler.getChapterContent)
	storyRouter.Get("/chapters/:id/session", handler.getUserSession)
	storyRouter.Post("/chapters/:id/start", handler.startSession)
	storyRouter.Post("/action", handler.submitAction)
}

func (h *storyHandler) getChapterList(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	resp, apiErr := h.uc.GetChapterList(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Daftar chapter berhasil dimuat", resp)
}

func (h *storyHandler) getChapterContent(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	chapterIDStr := ctx.Params("id")
	chapterID, err := uuid.Parse(chapterIDStr)
	if err != nil {
		return response.Error(ctx, response.NewParamValidationError("id", "uuid"), err)
	}

	resp, apiErr := h.uc.GetChapterContent(ctx.Context(), userID, chapterID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Konten chapter berhasil dimuat", resp)
}

func (h *storyHandler) getUserSession(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	chapterIDStr := ctx.Params("id")
	chapterID, err := uuid.Parse(chapterIDStr)
	if err != nil {
		return response.Error(ctx, response.NewParamValidationError("id", "uuid"), err)
	}

	resp, apiErr := h.uc.GetUserSession(ctx.Context(), userID, chapterID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Prorgesmu berhasil dipulihkan", resp)
}

func (h *storyHandler) startSession(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	chapterIDStr := ctx.Params("id")
	chapterID, err := uuid.Parse(chapterIDStr)
	if err != nil {
		return response.Error(ctx, response.NewParamValidationError("id", "uuid"), err)
	}

	if apiErr := h.uc.StartSession(ctx.Context(), userID, chapterID); apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Permainan dimulai, semangat ya!", nil)
}

func (h *storyHandler) submitAction(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	req := new(dto.StoryActionRequest)
	if err := ctx.BodyParser(&req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Data yang kamu kirim belum pas, coba cek lagi ya"), err)
	}

	if err := h.val.ValidateStruct(req); err != nil {
		return response.Error(ctx, response.NewValidationError(err), err)
	}

	resp, apiErr := h.uc.SubmitAction(ctx.Context(), userID, req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Aksimu berhasil diproses!", resp)
}
