package settingsruntime

import (
	"errors"
	"strings"
)

func CacheKey(entity, key string) string {
	return strings.ToUpper(strings.TrimSpace(entity)) + ":" + strings.TrimSpace(key)
}

func normalizeSpecKey(spec KeySpec) (string, error) {
	entity := strings.ToUpper(strings.TrimSpace(spec.Entity))
	key := strings.TrimSpace(spec.Key)
	if entity == "" || key == "" {
		return "", errors.New("entity and key are required")
	}
	return CacheKey(entity, key), nil
}
