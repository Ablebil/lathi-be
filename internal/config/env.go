package config

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv        string        `mapstructure:"APP_ENV"`
	AppHost       string        `mapstructure:"APP_HOST"`
	AppPort       int           `mapstructure:"APP_PORT"`
	DbHost        string        `mapstructure:"DB_HOST"`
	DbPort        int           `mapstructure:"DB_PORT"`
	DbName        string        `mapstructure:"DB_NAME"`
	DbUser        string        `mapstructure:"DB_USER"`
	DbPassword    string        `mapstructure:"DB_PASSWORD"`
	DbSSLMode     string        `mapstructure:"DB_SSLMODE"`
	RedisHost     string        `mapstructure:"REDIS_HOST"`
	RedisPort     int           `mapstructure:"REDIS_PORT"`
	RedisPassword string        `mapstructure:"REDIS_PASSWORD"`
	RedisDb       int           `mapstructure:"REDIS_DB"`
	AccessSecret  string        `mapstructure:"ACCESS_SECRET"`
	RefreshSecret string        `mapstructure:"REFRESH_SECRET"`
	AccessTtl     time.Duration `mapstructure:"ACCESS_TTL"`
	RefreshTtl    time.Duration `mapstructure:"REFRESH_TTL"`
	SmtpHost      string        `mapstructure:"SMTP_HOST"`
	SmtpPort      int           `mapstructure:"SMTP_PORT"`
	SmtpUsername  string        `mapstructure:"SMTP_USERNAME"`
	SmtpPassword  string        `mapstructure:"SMTP_PASSWORD"`
	VerifUrl      string        `mapstructure:"VERIF_URL"`
	VerifTokenTtl time.Duration `mapstructure:"VERIF_TOKEN_TTL"`
	FeUrl         string        `mapstructure:"FE_URL"`
}

func New() (*Env, error) {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile(".env")

	if err := v.ReadInConfig(); err != nil {
		slog.Error("failed to read config", "error", err)
		return nil, err
	}

	var env Env
	v.Unmarshal(&env)

	return &env, nil
}
