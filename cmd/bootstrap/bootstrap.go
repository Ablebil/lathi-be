package bootstrap

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/Ablebil/lathi-be/db/migration"
	"github.com/Ablebil/lathi-be/db/seed"
	"github.com/Ablebil/lathi-be/internal/config"
	cronJob "github.com/Ablebil/lathi-be/internal/infra/cron"
	"github.com/Ablebil/lathi-be/internal/infra/fiber"
	"github.com/Ablebil/lathi-be/internal/infra/minio"
	"github.com/Ablebil/lathi-be/internal/infra/postgresql"
	"github.com/Ablebil/lathi-be/internal/infra/redis"
	"github.com/Ablebil/lathi-be/internal/middleware"
	"github.com/Ablebil/lathi-be/pkg/bcrypt"
	"github.com/Ablebil/lathi-be/pkg/jwt"
	"github.com/Ablebil/lathi-be/pkg/mail"
	"github.com/Ablebil/lathi-be/pkg/validator"

	authHdl "github.com/Ablebil/lathi-be/internal/app/auth/handler"
	authUc "github.com/Ablebil/lathi-be/internal/app/auth/usecase"

	userHdl "github.com/Ablebil/lathi-be/internal/app/user/handler"
	userRepo "github.com/Ablebil/lathi-be/internal/app/user/repository"
	userUc "github.com/Ablebil/lathi-be/internal/app/user/usecase"

	storyHdl "github.com/Ablebil/lathi-be/internal/app/story/handler"
	storyRepo "github.com/Ablebil/lathi-be/internal/app/story/repository"
	storyUc "github.com/Ablebil/lathi-be/internal/app/story/usecase"

	dictHdl "github.com/Ablebil/lathi-be/internal/app/dictionary/handler"
	dictRepo "github.com/Ablebil/lathi-be/internal/app/dictionary/repository"
	dictUc "github.com/Ablebil/lathi-be/internal/app/dictionary/usecase"

	lbHdl "github.com/Ablebil/lathi-be/internal/app/leaderboard/handler"
	lbRepo "github.com/Ablebil/lathi-be/internal/app/leaderboard/repository"
	lbUc "github.com/Ablebil/lathi-be/internal/app/leaderboard/usecase"
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

	storage, err := minio.New(env)
	if err != nil {
		panic(err)
	}

	cache := redis.New(env)

	handleArgs(env)

	app := fiber.New(env)
	v1 := app.Group("/api/v1")

	val := validator.NewValidator()
	bcrypt := bcrypt.NewBcrypt()
	mail := mail.NewMail(env)
	jwt := jwt.NewJWT(env)
	mw := middleware.NewMiddleware(jwt)

	// auth module
	userRepository := userRepo.NewUserRepository(db)
	authUsecase := authUc.NewAuthUsecase(userRepository, bcrypt, mail, cache, jwt, env)
	authHdl.NewAuthHandler(v1, val, env, authUsecase)

	// leaderboard module
	leaderboardRepository := lbRepo.NewLeaderboardRepository(db, cache)
	leaderboardUsecase := lbUc.NewLeaderboardUsecase(leaderboardRepository, storage)
	lbHdl.NewLeaderboardHandler(v1, leaderboardUsecase)

	// story module
	storyRepository := storyRepo.NewStoryRepository(db)
	storyUsecase := storyUc.NewStoryUsecase(storyRepository, userRepository, leaderboardRepository, storage, env)
	storyHdl.NewStoryHandler(v1, val, mw, storyUsecase)

	// dictionary module
	dictionaryRepository := dictRepo.NewDictionaryRepository(db)
	dictionaryUsecase := dictUc.NewDictionaryUsecase(dictionaryRepository, env)
	dictHdl.NewDictionaryHandler(v1, val, mw, dictionaryUsecase)

	// user module
	userUsecase := userUc.NewUserUsecase(userRepository, storyRepository, dictionaryRepository, leaderboardRepository, storage, cache, env)
	userHdl.NewUserHandler(v1, val, mw, userUsecase)

	cron := cronJob.NewCronJob(userRepository, leaderboardRepository)
	cron.Start()

	return app.Listen(fmt.Sprintf("%s:%d", env.AppHost, env.AppPort))
}

func handleArgs(env *config.Env) {
	migrateCmd := flag.NewFlagSet("migrate", flag.ExitOnError)
	seedCmd := flag.NewFlagSet("seed", flag.ExitOnError)

	migrateAction := migrateCmd.String("action", "", "specify 'up' or 'down' for migration")
	seedDomain := seedCmd.String("domain", "", "specify a domain for seeding (optional)")

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
		case "seed":
			if err := seedCmd.Parse(os.Args[2:]); err != nil {
				slog.Error("unable to parse seed command", "error", err)
			}

			seed.Seed(env, *seedDomain)
			os.Exit(1)
		}
	}
}
