package util

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type Config struct {
	DBDriver            string        `mapstructure:"DB_DRIVER"`
	DBSource            string        `mapstructure:"DB_SOURCE"`
	ServerAddress       string        `mapstructure:"SERVER_ADDRESS"`
	TokenKey            string        `mapstructure:"TOKEN_KEY"`
	AccessTokenDuration time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
	GithubClientId      string        `mapstructure:"GITHUB_CLIENT_ID"`
	GithubClientSecret  string        `mapstructure:"GITHUB_CLIENT_SECRET"`
	RedirectUrl         string        `mapstructure:"REDIRECT_URL"`
	GinMode             string        `mapstructure:"GIN_MODE"`
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
	if config.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	return
}
