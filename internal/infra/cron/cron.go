package cron

import (
	"context"
	"log/slog"
	"time"

	"github.com/Ablebil/lathi-be/internal/domain/contract"
	"github.com/robfig/cron/v3"
)

type CronJobItf interface {
	Start()
	Stop()
}

type cronJob struct {
	userRepo contract.UserRepositoryItf
	lbRepo   contract.LeaderboardRepositoryItf
	cron     *cron.Cron
}

func NewCronJob(userRepo contract.UserRepositoryItf, lbRepo contract.LeaderboardRepositoryItf) CronJobItf {
	return &cronJob{
		userRepo: userRepo,
		lbRepo:   lbRepo,
		cron:     cron.New(),
	}
}

func (c *cronJob) Start() {
	slog.Info("starting cron jobs...")

	c.deleteUnverifiedUsersJob()
	c.rebuildLeaderboardJob()

	c.cron.Start()
	slog.Info("all cron jobs started successfully")
}

func (c *cronJob) Stop() {
	c.cron.Stop()
	slog.Info("all cron jobs stopped")
}

func (c *cronJob) deleteUnverifiedUsersJob() {
	// run every hour at minute 0
	_, err := c.cron.AddFunc("0 * * * *", func() {
		ctx := context.Background()
		threshold := time.Now().Add(-24 * time.Hour)

		deleted, err := c.userRepo.DeleteUnverifiedUsers(ctx, threshold)
		if err != nil {
			slog.Error("failed to delete unverified users", "error", err)
			return
		}

		if deleted > 0 {
			slog.Info("deleted unverified users", "count", deleted, "threshold", threshold)
		}
	})

	if err != nil {
		slog.Error("failed to register delete unverified users job", "error", err)
		return
	}

	slog.Info("delete unverified users job registered", "schedule", "every hour")
}

func (c *cronJob) rebuildLeaderboardJob() {
	// run every 6 hours
	_, err := c.cron.AddFunc("0 */6 * * *", func() {
		ctx := context.Background()

		if err := c.lbRepo.RebuildLeaderboard(ctx); err != nil {
			slog.Error("failed to rebuild leaderboard", "error", err)
			return
		}

		slog.Info("leaderboard rebuilt successfully")
	})

	if err != nil {
		slog.Error("failed to register rebuild leaderboard job", "error", err)
		return
	}

	slog.Info("rebuild leaderboard job registered", "schedule", "every 6 hours")
}
