package middleware

import (
	"github.com/Ablebil/lathi-be/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type middleware struct {
	jwt jwt.JwtItf
}

type MiddlewareItf interface {
	Authenticate(ctx *fiber.Ctx) error
}

func NewMiddleware(jwt jwt.JwtItf) *middleware {
	return &middleware{
		jwt: jwt,
	}
}
