package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

func getErrorMessageForEnv(key string, value string) string {
	return fmt.Sprintf("Invalid value for environment variable \"%s\": %s", key, value)
}

func GetEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	log.Fatalf("%s environment variable must be provided", key)

	return ""
}

func GetEnvIntOrDefault(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		log.Fatal(getErrorMessageForEnv(key, value))
	}

	return intValue
}

func GetEnvOrDefault(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	return value
}
