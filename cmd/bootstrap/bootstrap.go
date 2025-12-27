package bootstrap

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Ablebil/lathi-be/db/migration"
	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/internal/infra/fiber"
	"github.com/Ablebil/lathi-be/internal/infra/postgresql"
	"github.com/Ablebil/lathi-be/internal/infra/redis"
	"github.com/Ablebil/lathi-be/internal/middleware"
	"github.com/Ablebil/lathi-be/pkg/bcrypt"
	"github.com/Ablebil/lathi-be/pkg/jwt"
	"github.com/Ablebil/lathi-be/pkg/mail"
	"github.com/Ablebil/lathi-be/pkg/validator"

	authHdl "github.com/Ablebil/lathi-be/internal/app/auth/handler"
	authUc "github.com/Ablebil/lathi-be/internal/app/auth/usecase"

	userRepo "github.com/Ablebil/lathi-be/internal/app/user/repository"
)

func Start() error {
	env, err := config.New()
	if err != nil {
		panic(err)
	}

	db, err := postgresql.New(env)
	if err != nil {
		panic(err)
	}

	handleArgs(env)

	app := fiber.New(env)
	v1 := app.Group("/api/v1")

	cache := redis.New(env)
	val := validator.NewValidator()
	bcrypt := bcrypt.NewBcrypt()
	mail := mail.NewMail(env)
	jwt := jwt.NewJwt(env)
	mdw := middleware.NewMiddleware(jwt)

	// auth module
	userRepository := userRepo.NewUserRepository(db)
	authUsecase := authUc.NewAuthUsecase(userRepository, bcrypt, mail, cache, jwt, env)
	authHdl.NewAuthHandler(v1, val, authUsecase)

	return app.Listen(fmt.Sprintf("%s:%d", env.AppHost, env.AppPort))
}

func handleArgs(env *config.Env) {
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)

	migrateAction := migrateCmd.String("action", "", "specify 'up' or 'down' for migration")

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "migrate":
			if err := migrateCmd.Parse(os.Args[2:]); err != nil {
				slog.Error("unable to parse migrate command", "error", err)
			}

			if *migrateAction == "" {
				slog.Error("migration action is required")
			}

			migration.Migrate(env, *migrateAction)
			os.Exit(1)
		}
	}
}
