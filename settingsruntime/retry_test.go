package settingsruntime

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSanitizeRetry_Defaults(t *testing.T) {
	got := sanitizeRetry(RetryConfig{})
	if got.InitialInterval <= 0 || got.MaxInterval <= 0 || got.MaxElapsed <= 0 {
		t.Fatalf("expected positive defaults: %#v", got)
	}
	if got.JitterFactor <= 0 || got.JitterFactor > 1 {
		t.Fatalf("expected jitter in (0,1]: %#v", got)
	}
}

func TestNextBackoff(t *testing.T) {
	if got := nextBackoff(2*time.Second, 10*time.Second); got != 4*time.Second {
		t.Fatalf("unexpected backoff: %s", got)
	}
	if got := nextBackoff(8*time.Second, 10*time.Second); got != 10*time.Second {
		t.Fatalf("expected max cap: %s", got)
	}
}

func TestAddJitter(t *testing.T) {
	base := 100 * time.Millisecond
	for range 20 {
		got := addJitter(base, 0.2)
		if got < time.Millisecond {
			t.Fatalf("jitter below minimum: %s", got)
		}
	}
}

func TestWaitForCoreReady_SuccessAfterRetry(t *testing.T) {
	attempts := 0
	err := waitForCoreReady(
		context.Background(),
		func(context.Context) error {
			attempts++
			if attempts < 3 {
				return errors.New("not ready")
			}
			return nil
		},
		RetryConfig{
			InitialInterval: time.Millisecond,
			MaxInterval:     time.Millisecond,
			MaxElapsed:      50 * time.Millisecond,
			JitterFactor:    0,
		},
	)
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}
}
