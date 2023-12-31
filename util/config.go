package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Environtment         string        `mapstructure:"ENVIRONTMENT"`
	RedisAddress         string        `mapstructure:"REDIS_ADDR"`
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	HTTPServerAddress    string        `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress    string        `mapstructure:"GRPC_SERVER_ADDRESS"`
	TokenKey             string        `mapstructure:"TOKEN_KEY"`
	AccessTokenDuration  time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	GmailUsername        string        `mapstructure:"GMAIL_USERNAME"`
	GmailPassword        string        `mapstructure:"GMAIL_PASSWORD"`
	GmailName            string        `mapstructure:"GMAIL_NAME"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	// viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()

	if err != nil {
		return
	}

	// opt := viper.DecodeHook(mapstructure.StringToTimeDurationHookFunc())

	err = viper.Unmarshal(&config)
	return
}
