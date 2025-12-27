package middleware

import (
	"strings"

	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/gofiber/fiber/v2"
)

func (m *middleware) Authenticate(ctx *fiber.Ctx) error {
	header := ctx.Get("Authorization")
	if header == "" {
		return response.Error(ctx, response.ErrUnauthorized("you're not logged in. please login to continue"), nil)
	}

	parts := strings.Fields(header)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return response.Error(ctx, response.ErrUnauthorized("invalid authorization header format"), nil)
	}

	token := parts[1]

	validate, err := m.jwt.ParseAccessToken(token)
	if err != nil {
		return response.Error(ctx, response.ErrUnauthorized("invalid or expired token"), nil)
	}

	ctx.Locals("user_id", validate.Subject)
	ctx.Locals("username", validate.Username)
	ctx.Locals("email", validate.Email)

	return ctx.Next()
}
