package utils

import (
	"strconv"
)

// ParseIntStr parses a string into an int64 number.
func ParseIntStr(intStr string) (int64, error) {
	number, err := strconv.ParseInt(intStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return number, nil
}
