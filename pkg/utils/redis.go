package utils

import "fmt"

func GenerateRedisKey(basePrefix string, uniqueKey string) string {
	return fmt.Sprintf("%s:%s", basePrefix, uniqueKey)
}

func GetRedisKey(basePrefix string, key string) string {
	return fmt.Sprintf("%s:%s", basePrefix, key)
}
