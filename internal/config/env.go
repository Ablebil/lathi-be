package config

import (
	"log/slog"

	"github.com/spf13/viper"
)

type Env struct {
	AppEnv     string `mapstructure:"APP_ENV"`
	AppHost    string `mapstructure:"APP_HOST"`
	AppPort    int    `mapstructure:"APP_PORT"`
	DbHost     string `mapstructure:"DB_HOST"`
	DbPort     int    `mapstructure:"DB_PORT"`
	DbName     string `mapstructure:"DB_NAME"`
	DbUser     string `mapstructure:"DB_USER"`
	DbPassword string `mapstructure:"DB_PASSWORD"`
	DbSSLMode  string `mapstructure:"DB_SSLMODE"`
	FeUrl      string `mapstructure:"FE_URL"`
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
