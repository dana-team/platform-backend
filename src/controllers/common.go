package controllers

import "github.com/dana-team/platform-backend/src/types"

func convertKeyValueToMap(kvList []types.KeyValue) map[string]string {
	values := make(map[string]string)
	for _, kv := range kvList {
		values[kv.Key] = kv.Value
	}
	return values
}

func convertMapToKeyValue(values map[string]string) []types.KeyValue {
	var kvList []types.KeyValue
	for k, v := range values {
		kvList = append(kvList, types.KeyValue{Key: k, Value: v})
	}
	return kvList
}
