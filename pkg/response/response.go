package response

import (
	"log/slog"

	"github.com/gofiber/fiber/v2"
)

func Success(ctx *fiber.Ctx, status int, data any) error {
	return ctx.Status(status).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}

func Error(ctx *fiber.Ctx, apiErr *APIError, err error) error {
	slog.Error("API error",
		"status", apiErr.Status,
		"type", apiErr.Type,
		"message", apiErr.Message,
		"detail", apiErr.Detail,
		"fields", apiErr.Fields,
		"error", err,
		"path", ctx.Path(),
	)

	resp := fiber.Map{
		"success": false,
		"error": fiber.Map{
			"type":    apiErr.Type,
			"message": apiErr.Message,
			"detail":  apiErr.Detail,
		},
	}

	if len(apiErr.Fields) > 0 {
		resp["error"].(fiber.Map)["fields"] = apiErr.Fields
	}

	return ctx.Status(apiErr.Status).JSON(resp)
}
