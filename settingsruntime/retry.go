package settingsruntime

import (
	"context"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

func sanitizeRetry(cfg RetryConfig) RetryConfig {
	if cfg.InitialInterval <= 0 {
		cfg.InitialInterval = 1 * time.Second
	}
	if cfg.MaxInterval <= 0 {
		cfg.MaxInterval = 30 * time.Second
	}
	if cfg.MaxElapsed <= 0 {
		cfg.MaxElapsed = 2 * time.Minute
	}
	if cfg.JitterFactor <= 0 {
		cfg.JitterFactor = 0.2
	}
	if cfg.JitterFactor > 1 {
		cfg.JitterFactor = 1
	}
	return cfg
}

func waitForCoreReady(ctx context.Context, healthCheck func(context.Context) error, retry RetryConfig) error {
	start := time.Now()
	interval := retry.InitialInterval

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		if err := healthCheck(ctx); err == nil {
			return nil
		}

		if time.Since(start) >= retry.MaxElapsed {
			return fmt.Errorf("timeout after %s", retry.MaxElapsed)
		}

		sleep := addJitter(interval, retry.JitterFactor)
		timer := time.NewTimer(sleep)
		select {
		case <-ctx.Done():
			timer.Stop()
			return ctx.Err()
		case <-timer.C:
		}

		interval = nextBackoff(interval, retry.MaxInterval)
	}
}

func nextBackoff(current, max time.Duration) time.Duration {
	next := current * 2
	if next > max {
		return max
	}
	return next
}

func addJitter(base time.Duration, factor float64) time.Duration {
	if base <= 0 || factor <= 0 {
		return base
	}

	spread := float64(base) * factor
	delta := (randomUnitFloat64()*2 - 1) * spread
	jittered := float64(base) + delta
	if jittered < float64(time.Millisecond) {
		jittered = float64(time.Millisecond)
	}
	if jittered > math.MaxInt64 {
		jittered = math.MaxInt64
	}
	return time.Duration(jittered)
}

func randomUnitFloat64() float64 {
	var buf [8]byte
	if _, err := rand.Read(buf[:]); err != nil {
		// Fallback deterministic value in the unlikely case CSPRNG fails.
		return 0.5
	}
	value := binary.BigEndian.Uint64(buf[:])
	return float64(value) / float64(^uint64(0))
}
