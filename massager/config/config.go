package config

import (
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	SecretToken *SecretToken `mapstructure:"SECRET_TOKEN"`
	SwaggerUser string       `mapstructure:"SWAGGER_USER"`
	SwaggerPass string       `mapstructure:"SWAGGER_PASS"`
	PgConfig    *PgConfig    `mapstructure:"PG_CONFIG"`
}

type SecretToken struct {
	TokenSymmetricKey        string        `mapstructure:"JWT_SECRET_KEY"`
	RefreshTokenSymmetricKey string        `mapstructure:"JWT_REFRESH_SECRET_KEY"`
	AccessTokenDuration      time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	RefreshTokenDuration     time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
}

type PgConfig struct {
	Host     string `mapstructure:"POSTGRES_HOST"`
	User     string `mapstructure:"POSTGRES_USER"`
	Password string `mapstructure:"POSTGRES_PASSWORD"`
	Port     int    `mapstructure:"POSTGRES_PORT"`
	Database string `mapstructure:"POSTGRES_DATABASE"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	swaggerUser := config.SwaggerUser
	swaggerPass := config.SwaggerPass
	if swaggerUser == "" || swaggerPass == "" {
		log.Fatal("Environment variables SWAGGER_USER and SWAGGER_PASS must be set")
	}

	return
}
