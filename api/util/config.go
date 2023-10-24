package util

import (
	"time"

	"github.com/spf13/viper"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver             string        `mapstructure:"DB_DRIVER"`
	DBSource             string        `mapstructure:"DB_SOURCE"`
	ServerAddress        string        `mapstructure:"SERVER_ADDRESS"`
	SessionTokenDuration time.Duration `mapstructure:"SESSION_TOKEN_DURATION"`
	RefreshTokenDuration time.Duration `mapstructure:"REFRESH_TOKEN_DURATION"`
	TestAccountUsername1 string        `mapstructure:"TEST_ACCOUNT_USERNAME_1"`
	TestAccountEmail1    string        `mapstructure:"TEST_ACCOUNT_EMAIL_1"`
	TestAccountUsername2 string        `mapstructure:"TEST_ACCOUNT_USERNAME_2"`
	TestAccountEmail2    string        `mapstructure:"TEST_ACCOUNT_EMAIL_2"`
	TestAccountUsername3 string        `mapstructure:"TEST_ACCOUNT_USERNAME_3"`
	TestAccountEmail3    string        `mapstructure:"TEST_ACCOUNT_EMAIL_3"`
	TestAccountPassword  string        `mapstructure:"TEST_ACCOUNT_PASSWORD"`
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
