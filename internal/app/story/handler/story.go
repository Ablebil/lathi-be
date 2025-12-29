package handler

import (
	"github.com/Ablebil/lathi-be/internal/domain/contract"
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
}

func (h *storyHandler) getChapterList(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("invalid user session"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	resp, apiErr := h.uc.GetChapterList(ctx.Context(), userID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "chapter list retrieved successfully", resp)
}

func (h *storyHandler) getChapterContent(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("invalid user session"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	chapterIDStr := ctx.Params("id")
	chapterID, err := uuid.Parse(chapterIDStr)
	if err != nil {
		return response.Error(ctx, response.ErrBadRequest("invalid chapter ID"), err)
	}

	resp, apiErr := h.uc.GetChapterContent(ctx.Context(), userID, chapterID)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "chapter content retrieved successfully", resp)
}
