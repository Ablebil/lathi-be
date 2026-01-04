package fiber

import (
	"time"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func New(env *config.Env) *fiber.App {
	app := fiber.New(fiber.Config{
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  5 * time.Second,
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}

			return ctx.Status(code).JSON(fiber.Map{
				"success": false,
				"error": fiber.Map{
					"type":    "app_error",
					"message": err.Error(),
					"status":  code,
				},
			})
		},
	})

	app.Use(recover.New())
	app.Use(requestid.New())
	app.Use(helmet.New())
	app.Use(healthcheck.New())
	app.Use(compress.New(compress.Config{Level: compress.LevelBestSpeed}))
	app.Use(cors.New(cors.Config{
		AllowOrigins: func() string {
			if env.AppEnv == "development" {
				return "*"
			}
			return env.FEURL
		}(),
		AllowMethods:     "GET,POST,PUT,PATCH,DELETE,OPTIONS",
		AllowHeaders:     "Content-Type, Authorization",
		AllowCredentials: env.AppEnv == "production",
	}))
	app.Use(logger.New(logger.Config{
		Format: "${time} | ${status} | ${method} | ${path} | ${latency}\n",
	}))

	return app
}
