package utils

import (
	"github.com/dana-team/platform-backend/internal/types"
)

// ConvertKeyValueToMap converts a slice of KeyValue pairs to a map
// with string keys and values.
func ConvertKeyValueToMap(kvList []types.KeyValue) map[string]string {
	values := make(map[string]string)
	for _, kv := range kvList {
		values[kv.Key] = kv.Value
	}
	return values
}

// ConvertMapToKeyValue converts a map with string keys and values
// to a slice of KeyValue pairs.
func ConvertMapToKeyValue(values map[string]string) []types.KeyValue {
	var kvList []types.KeyValue
	for k, v := range values {
		kvList = append(kvList, types.KeyValue{Key: k, Value: v})
	}
	return kvList
}
