package util

import (
	"fmt"
	"os"
	"time"
)

// Config stores all configuration of the application.
// The values are read by viper from a config file or environment variables.
type Config struct {
	DBDriver             string
	DBSource             string
	ServerAddress        string
	SessionTokenDuration time.Duration
	TestAccountUsername  string
	TestAccountEmail     string
	TestAccountPassword  string
}

func getEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s environment variable not set", key)
	}
	return value, nil
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig() (config Config, err error) {
	dbDriver, err := getEnv("DB_DRIVER")
	if err != nil {
		return
	}

	dbSource, err := getEnv("DB_SOURCE")
	if err != nil {
		return
	}

	serverAddress, err := getEnv("SERVER_ADDRESS")
	if err != nil {
		return
	}

	sessionTokenDurationStr, err := getEnv("SESSION_TOKEN_DURATION")
	if err != nil {
		return
	}

	sessionTokenDuration, err := time.ParseDuration(sessionTokenDurationStr)
	if err != nil {
		return
	}

	testAccountUsername, err := getEnv("TEST_ACCOUNT_USERNAME")
	if err != nil {
		return
	}

	testAccountEmail, err := getEnv("TEST_ACCOUNT_EMAIL")
	if err != nil {
		return
	}

	testAccountPassword, err := getEnv("TEST_ACCOUNT_PASSWORD")
	if err != nil {
		return
	}

	config = Config{
		DBDriver:             dbDriver,
		DBSource:             dbSource,
		ServerAddress:        serverAddress,
		SessionTokenDuration: sessionTokenDuration,
		TestAccountUsername:  testAccountUsername,
		TestAccountEmail:     testAccountEmail,
		TestAccountPassword:  testAccountPassword,
	}

	return
}
