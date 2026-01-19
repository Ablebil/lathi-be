package handler

import (
	"time"

	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/Ablebil/lathi-be/internal/domain/dto"
	"github.com/Ablebil/lathi-be/internal/middleware"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/Ablebil/lathi-be/pkg/validator"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type dictionaryHandler struct {
	val validator.ValidatorItf
	uc  contract.DictionaryUsecaseItf
}

func NewDictionaryHandler(router fiber.Router, validator validator.ValidatorItf, mw middleware.MiddlewareItf, dicttionaryUc contract.DictionaryUsecaseItf) {
	handler := dictionaryHandler{
		val: validator,
		uc:  dicttionaryUc,
	}

	dictionaryRouter := router.Group("/dictionaries", mw.Authenticate)
	dictionaryRouter.Get("/", mw.RateLimit(60, 1*time.Minute, "dict_list"), handler.getDictionaryList)
}

func (h *dictionaryHandler) getDictionaryList(ctx *fiber.Ctx) error {
	userIDStr, ok := ctx.Locals("user_id").(string)
	if !ok {
		return response.Error(ctx, response.ErrUnauthorized("Kamu belum login, yuk login dulu"), nil)
	}
	userID, _ := uuid.Parse(userIDStr)

	allowedParams := map[string]bool{
		"search": true,
		"page":   true,
		"limit":  true,
	}

	queryParams := ctx.Queries()
	for k, v := range queryParams {
		if !allowedParams[k] {
			return response.Error(ctx, response.ErrBadRequest("Query parameter '"+k+"' ga dikenali"), nil)
		}

		if v == "" {
			return response.Error(ctx, response.ErrBadRequest("Query parameter '"+k+"' ga boleh kosong"), nil)
		}
	}

	req := new(dto.DictionaryListRequest)
	if err := ctx.QueryParser(req); err != nil {
		return response.Error(ctx, response.ErrBadRequest("Format query ga valid"), err)
	}

	resp, apiErr := h.uc.GetDictionaryList(ctx.Context(), userID, req)
	if apiErr != nil {
		return response.Error(ctx, apiErr, nil)
	}

	return response.Success(ctx, fiber.StatusOK, "Kamus berhasil dimuat", resp)
}
