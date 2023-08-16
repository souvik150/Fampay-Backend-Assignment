package initializers

import (
	"github.com/spf13/viper"
	"time"
)

type Config struct {
	Port               string        `mapstructure:"PORT"`
	DBHost             string        `mapstructure:"POSTGRES_HOST"`
	DBUserName         string        `mapstructure:"POSTGRES_USER"`
	DBUserPassword     string        `mapstructure:"POSTGRES_PASSWORD"`
	DBName             string        `mapstructure:"POSTGRES_DB"`
	DBPort             string        `mapstructure:"POSTGRES_PORT"`
	AccessTokenSecret  string        `mapstructure:"ACCESS_TOKEN_SECRET"`
	RefreshTokenSecret string        `mapstructure:"REFRESH_TOKEN_SECRET"`
	AccessTokenExpiry  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRY"`
	RefreshTokenExpiry time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRY"`
	ClientOrigin       string        `mapstructure:"CLIENT_ORIGIN"`
	Email              string        `mapstructure:"EMAIL_ID"`
	EmailPassword      string        `mapstructure:"EMAIL_PASSWORD"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		return
	}

	config.AccessTokenExpiry, err = time.ParseDuration(config.AccessTokenExpiry.String())
	if err != nil {
		return
	}

	config.RefreshTokenExpiry, err = time.ParseDuration(config.RefreshTokenExpiry.String())
	if err != nil {
		return
	}
	return
}
