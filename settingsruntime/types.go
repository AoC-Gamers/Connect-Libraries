package settingsruntime

import (
	"context"
	"time"
)

const defaultHTTPTimeout = 5 * time.Second

type ValueValidator func(string) error

type KeySpec struct {
	Entity   string
	Key      string
	Validate ValueValidator
}

type RetryConfig struct {
	InitialInterval time.Duration
	MaxInterval     time.Duration
	MaxElapsed      time.Duration
	JitterFactor    float64
}

type BootstrapConfig struct {
	ServiceName string
	Retry       RetryConfig
	Required    []KeySpec
	HealthCheck func(context.Context) error
	Getter      func(context.Context, string, string) (string, error)
}
