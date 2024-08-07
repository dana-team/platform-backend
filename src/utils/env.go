package utils

import (
	"fmt"
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
		return false, fmt.Errorf("failed to parse %q", valStr)
	}

	return valBool, nil
}

// GetEnvNumber retrieves the value of the environment variable named by the key.
// If the variable is empty or not set, it returns the default value.
func GetEnvNumber(key string, defaultValue int) (int, error) {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultValue, nil
	}

	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		return 0, fmt.Errorf("failed to parse %q as integer", valStr)
	}

	return valInt, nil
}
