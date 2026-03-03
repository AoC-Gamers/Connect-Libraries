package settingsruntime

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestBootstrap_Success(t *testing.T) {
	ctx := context.Background()
	values, err := Bootstrap(ctx, BootstrapConfig{
		ServiceName: "svc",
		Retry: RetryConfig{
			InitialInterval: time.Millisecond,
			MaxInterval:     time.Millisecond,
			MaxElapsed:      20 * time.Millisecond,
			JitterFactor:    0,
		},
		Required: []KeySpec{
			{Entity: "CONFIG", Key: "feature.enabled", Validate: BoolValidator()},
		},
		HealthCheck: func(context.Context) error { return nil },
		Getter: func(context.Context, string, string) (string, error) {
			return "true", nil
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if values["CONFIG:feature.enabled"] != "true" {
		t.Fatalf("unexpected value map: %#v", values)
	}
}

func TestBootstrap_DuplicateKey(t *testing.T) {
	ctx := context.Background()
	_, err := Bootstrap(ctx, BootstrapConfig{
		ServiceName: "svc",
		Required: []KeySpec{
			{Entity: "CONFIG", Key: "a"},
			{Entity: "config", Key: "a"},
		},
		HealthCheck: func(context.Context) error { return nil },
		Getter:      func(context.Context, string, string) (string, error) { return "x", nil },
	})
	if err == nil {
		t.Fatal("expected duplicate key error")
	}
}

func TestBootstrap_HealthFailsUntilTimeout(t *testing.T) {
	ctx := context.Background()
	_, err := Bootstrap(ctx, BootstrapConfig{
		ServiceName: "svc",
		Retry: RetryConfig{
			InitialInterval: time.Millisecond,
			MaxInterval:     time.Millisecond,
			MaxElapsed:      10 * time.Millisecond,
			JitterFactor:    0,
		},
		Required:    []KeySpec{},
		HealthCheck: func(context.Context) error { return errors.New("down") },
		Getter:      func(context.Context, string, string) (string, error) { return "", nil },
	})
	if err == nil {
		t.Fatal("expected readiness timeout error")
	}
}

func TestBootstrap_InvalidRequiredValue(t *testing.T) {
	ctx := context.Background()
	_, err := Bootstrap(ctx, BootstrapConfig{
		ServiceName: "svc",
		Required: []KeySpec{
			{Entity: "CONFIG", Key: "x", Validate: BoolValidator()},
		},
		HealthCheck: func(context.Context) error { return nil },
		Getter:      func(context.Context, string, string) (string, error) { return "nope", nil },
	})
	if err == nil {
		t.Fatal("expected validation error")
	}
}
