package settingsruntime

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

func Bootstrap(ctx context.Context, cfg BootstrapConfig) (map[string]string, error) {
	if cfg.HealthCheck == nil {
		return nil, errors.New("health check function is required")
	}
	if cfg.Getter == nil {
		return nil, errors.New("setting getter function is required")
	}

	retry := sanitizeRetry(cfg.Retry)
	if err := waitForCoreReady(ctx, cfg.HealthCheck, retry); err != nil {
		return nil, fmt.Errorf("%s: core readiness failed: %w", cfg.ServiceName, err)
	}

	values := make(map[string]string, len(cfg.Required))
	seen := make(map[string]struct{}, len(cfg.Required))

	for _, key := range cfg.Required {
		cacheKey, err := normalizeSpecKey(key)
		if err != nil {
			return nil, fmt.Errorf("%s: invalid key spec: %w", cfg.ServiceName, err)
		}

		if _, exists := seen[cacheKey]; exists {
			return nil, fmt.Errorf("%s: duplicated required setting %s", cfg.ServiceName, cacheKey)
		}
		seen[cacheKey] = struct{}{}

		value, err := cfg.Getter(ctx, key.Entity, key.Key)
		if err != nil {
			return nil, fmt.Errorf("%s: failed loading required setting %s: %w", cfg.ServiceName, cacheKey, err)
		}

		if strings.TrimSpace(value) == "" {
			return nil, fmt.Errorf("%s: required setting %s is empty", cfg.ServiceName, cacheKey)
		}

		if key.Validate != nil {
			if err := key.Validate(value); err != nil {
				return nil, fmt.Errorf("%s: invalid required setting %s: %w", cfg.ServiceName, cacheKey, err)
			}
		}

		values[cacheKey] = value
	}

	return values, nil
}
