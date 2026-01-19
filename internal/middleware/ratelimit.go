package middleware

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/pkg/response"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (m *middleware) RateLimit(limit int, window time.Duration, keyPrefix string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		identifier := getRateLimitIdentifier(ctx)
		if identifier == "" {
			identifier = ensureAnonID(ctx, m.env)
		}

		key := fmt.Sprintf("rl:%s:%s", keyPrefix, identifier)
		count, err := m.cache.Incr(ctx.Context(), key, window)
		if err != nil {
			return response.Error(ctx, response.ErrInternal("Coba lagi nanti ya!"), err)
		}

		ctx.Set("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
		ctx.Set("X-RateLimit-Remaining", fmt.Sprintf("%d", max(0, int64(limit)-count)))
		ctx.Set("X-RateLimit-Reset", fmt.Sprintf("%d", time.Now().Add(window).Unix()))

		if count > int64(limit) {
			if count == int64(limit)+1 {
				slog.Warn("rate limit exceeded",
					"key", keyPrefix,
					"identifier", identifier,
					"ip", ctx.IP())
			}

			ctx.Set("Retry-After", fmt.Sprintf("%d", int(window.Seconds())))
			return response.Error(ctx, response.ErrTooManyRequests("Terlalu banyak permintaan, coba lagi nanti ya"), nil)
		}
		return ctx.Next()
	}
}

func getRateLimitIdentifier(ctx *fiber.Ctx) string {
	if userID, ok := ctx.Locals("user_id").(string); ok && userID != "" {
		return userID
	}
	anonID := ctx.Cookies("anon_id")
	if anonID != "" {
		return anonID
	}

	ua := ctx.Get("User-Agent")
	lang := ctx.Get("Accept-Language")
	encoding := ctx.Get("Accept-Encoding")

	if ua == "" && lang == "" && encoding == "" {
		return ctx.IP()
	}

	return fmt.Sprintf("%s|%s|%s", ua, lang, encoding)
}

func ensureAnonID(ctx *fiber.Ctx, env *config.Env) string {
	anonID := ctx.Cookies("anon_id")
	if anonID == "" {
		anonID = uuid.NewString()
		ctx.Cookie(&fiber.Cookie{
			Name:     "anon_id",
			Value:    anonID,
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 30,
			HTTPOnly: true,
			Secure:   true,
			SameSite: func() string {
				if env.AppEnv == "production" {
					return fiber.CookieSameSiteLaxMode
				}
				return fiber.CookieSameSiteNoneMode
			}(),
		})
	}
	return anonID
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
