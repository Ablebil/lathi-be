package config

import (
	"log/slog"
	"time"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv           string        `mapstructure:"APP_ENV"`
	AppHost          string        `mapstructure:"APP_HOST"`
	AppPort          int           `mapstructure:"APP_PORT"`
	DBHost           string        `mapstructure:"DB_HOST"`
	DBPort           int           `mapstructure:"DB_PORT"`
	DBName           string        `mapstructure:"DB_NAME"`
	DBUser           string        `mapstructure:"DB_USER"`
	DBPassword       string        `mapstructure:"DB_PASSWORD"`
	DBSSLMode        string        `mapstructure:"DB_SSLMODE"`
	RedisHost        string        `mapstructure:"REDIS_HOST"`
	RedisPort        int           `mapstructure:"REDIS_PORT"`
	RedisPassword    string        `mapstructure:"REDIS_PASSWORD"`
	RedisDB          int           `mapstructure:"REDIS_DB"`
	AccessSecret     string        `mapstructure:"ACCESS_SECRET"`
	RefreshSecret    string        `mapstructure:"REFRESH_SECRET"`
	AccessTTL        time.Duration `mapstructure:"ACCESS_TTL"`
	RefreshTTL       time.Duration `mapstructure:"REFRESH_TTL"`
	SMTPHost         string        `mapstructure:"SMTP_HOST"`
	SMTPPort         int           `mapstructure:"SMTP_PORT"`
	SMTPUsername     string        `mapstructure:"SMTP_USERNAME"`
	SMTPPassword     string        `mapstructure:"SMTP_PASSWORD"`
	VerifURL         string        `mapstructure:"VERIF_URL"`
	VerifTokenTTL    time.Duration `mapstructure:"VERIF_TOKEN_TTL"`
	FEURL            string        `mapstructure:"FE_URL"`
	StorageEndpoint  string        `mapstructure:"STORAGE_ENDPOINT"`
	StoragePublicURL string        `mapstructure:"STORAGE_PUBLIC_URL"`
	StorageAccessKey string        `mapstructure:"STORAGE_ACCESS_KEY"`
	StorageSecretKey string        `mapstructure:"STORAGE_SECRET_KEY"`
	StorageBucket    string        `mapstructure:"STORAGE_BUCKET"`
	DefaultPageLimit int           `mapstructure:"DEFAULT_PAGE_LIMIT"`
	MaxPageLimit     int           `mapstructure:"MAX_PAGE_LIMIT"`
	DefaultAvatarURL string        `mapstructure:"DEFAULT_AVATAR_URL"`
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
