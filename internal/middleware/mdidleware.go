package middleware

import (
	"time"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/infra/redis"
	"github.com/Ablebil/lathi-be/pkg/jwt"
	"github.com/gofiber/fiber/v2"
)

type middleware struct {
	jwt   jwt.JWTItf
	cache redis.RedisItf
	env   *config.Env
}

type MiddlewareItf interface {
	Authenticate(ctx *fiber.Ctx) error
	RateLimit(limit int, window time.Duration, keyPrefix string) fiber.Handler
}

func NewMiddleware(jwt jwt.JWTItf, cache redis.RedisItf, env *config.Env) *middleware {
	return &middleware{
		jwt:   jwt,
		cache: cache,
		env:   env,
	}
}
