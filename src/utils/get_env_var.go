package utils

import (
	"os"
	"strconv"
)

// GetEnvBool retrieves the value of the environment variable named by the key.
// If the variable is empty or not set, it returns the default value.
func GetEnvBool(key string, defaultValue bool) (bool, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultValue, nil
	}

	valBool, err := strconv.ParseBool(valStr)
	if err != nil {
		return false, err
	}

	return valBool, nil
}
