package utils

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvOrDefault[T string | bool | int](key string, defaultValue T) T {
	envValue := os.Getenv(key)
	if envValue == "" {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case string:
		return any(envValue).(T)
	case int:
		intValue, err := strconv.Atoi(envValue)
		if err != nil {
			return defaultValue
		}
		return any(intValue).(T)
	case bool:
		boolValue, err := strconv.ParseBool(envValue)
		if err != nil {
			return defaultValue
		}
		return any(boolValue).(T)
	default:
		return defaultValue
	}
}

func GetEnv(key string) (string, error) {
	envValue := os.Getenv(key)
	if envValue == "" {
		return "", fmt.Errorf("environment variable doesnt exist")
	}
	return envValue, nil
}
