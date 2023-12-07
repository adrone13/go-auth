package config

import (
	"auth/internal/logger"
	"fmt"
	"os"
	"strconv"
)

func GetString(k string) string {
	v := os.Getenv(k)
	if v == "" {
		logger.Error(fmt.Sprintf(`ConfigService error: "%s" is not set up`, k))

		panic("config validation failed")
	}

	return v
}

func GetInt(k string) int {
	v := os.Getenv(k)
	if v == "" {
		logger.Error(fmt.Sprintf(`ConfigService error: "%s" is not set up`, k))

		panic("Config validation: value not set")
	}

	converted, err := strconv.Atoi(v)
	if err != nil {
		logger.Error(fmt.Sprintf(`ConfigService error: "%s" failed to convert to int`, k))

		panic("Config validation: failed to convert to int")
	}

	return converted
}
