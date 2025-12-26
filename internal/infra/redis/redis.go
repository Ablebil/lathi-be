package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/Ablebil/lathi-be/internal/config"
	r "github.com/redis/go-redis/v9"
)

type RedisItf interface {
	Set(ctx context.Context, key string, val any, exp time.Duration) error
	Get(ctx context.Context, key string, val any) error
	Del(ctx context.Context, key string) error
	Close() error
}

type redis struct {
	client *r.Client
}

func New(env *config.Env) RedisItf {
	client := r.NewClient(&r.Options{
		Addr:     fmt.Sprintf("%s:%d", env.RedisHost, env.RedisPort),
		Password: env.RedisPassword,
		DB:       env.RedisDb,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if _, err := client.Ping(ctx).Result(); err != nil {
		slog.Error("failed to connect to redis", "error", err)
		panic(err)
	}

	return &redis{client: client}
}

func (rd *redis) Set(ctx context.Context, key string, val any, exp time.Duration) error {
	return rd.client.Set(ctx, key, val, exp).Err()
}

func (rd *redis) Get(ctx context.Context, key string, val any) error {
	err := rd.client.Get(ctx, key).Scan(val)
	if err != nil {
		if errors.Is(err, r.Nil) {
			return fmt.Errorf("not found: %w", err)
		}
		return err
	}
	return nil
}

func (rd *redis) Del(ctx context.Context, key string) error {
	return rd.client.Del(ctx, key).Err()
}

func (rd *redis) Close() error {
	return rd.client.Close()
}
