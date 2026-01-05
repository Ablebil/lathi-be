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
	cron     *cron.Cron
}

func NewCronJob(userRepo contract.UserRepositoryItf) CronJobItf {
	return &cronJob{
		userRepo: userRepo,
		cron:     cron.New(),
	}
}

func (c *cronJob) Start() {
	slog.Info("starting cron jobs...")

	c.deleteUnverifiedUsersJob()

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
