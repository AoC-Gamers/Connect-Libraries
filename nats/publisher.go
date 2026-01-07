package connectnats

import (
	"encoding/json"
	"fmt"
	"time"

	natsio "github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
)

// Publisher provides methods to publish messages to NATS using both JetStream and Core
type Publisher struct {
	conn        *natsio.Conn
	js          natsio.JetStreamContext
	log         zerolog.Logger
	serviceName string
}

// PublisherOption configures a Publisher
type PublisherOption func(*Publisher)

// WithServiceName sets the service name for event metadata
func WithServiceName(name string) PublisherOption {
	return func(p *Publisher) {
		p.serviceName = name
	}
}

// NewPublisher creates a new Publisher from a NATS connection
func NewPublisher(conn *natsio.Conn, log zerolog.Logger, opts ...PublisherOption) (*Publisher, error) {
	if conn == nil {
		return nil, fmt.Errorf("NATS connection is nil")
	}

	js, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("failed to get JetStream context: %w", err)
	}

	p := &Publisher{
		conn:        conn,
		js:          js,
		log:         log.With().Str("component", "nats_publisher").Logger(),
		serviceName: "connect-service",
	}

	for _, opt := range opts {
		opt(p)
	}

	p.log.Debug().Str("service", p.serviceName).Msg("NATS Publisher initialized")
	return p, nil
}

// Event is the standard envelope for events published to JetStream
type Event struct {
	Type    string      `json:"type"`
	Version int         `json:"version"`
	Data    interface{} `json:"data"`
	Meta    EventMeta   `json:"meta"`
}

// EventMeta contains metadata about the event
type EventMeta struct {
	Timestamp int64  `json:"ts"`
	Service   string `json:"by"`
}

// PublishJetStream publishes a message to JetStream with ACK and persistence
// Use this for critical events that need guaranteed delivery
func (p *Publisher) PublishJetStream(subject string, data interface{}) error {
	if p.js == nil {
		return fmt.Errorf("JetStream not initialized")
	}

	b, err := json.Marshal(data)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("subject", subject).
			Msg("Failed to marshal data for JetStream")
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	p.log.Debug().
		Str("subject", subject).
		Int("bytes", len(b)).
		Msg("Publishing message to JetStream")

	_, err = p.js.Publish(subject, b)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("subject", subject).
			Msg("Failed to publish message to JetStream")
		return fmt.Errorf("failed to publish to JetStream: %w", err)
	}

	p.log.Debug().
		Str("subject", subject).
		Msg("Message published successfully to JetStream")

	return nil
}

// PublishCore publishes a message using NATS Core (fire-and-forget, no persistence)
// Use this for ephemeral events where occasional loss is acceptable
func (p *Publisher) PublishCore(subject string, data interface{}) error {
	if p.conn == nil {
		return fmt.Errorf("NATS connection not initialized")
	}

	b, err := json.Marshal(data)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("subject", subject).
			Msg("Failed to marshal data for NATS Core")
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	p.log.Debug().
		Str("subject", subject).
		Int("bytes", len(b)).
		Msg("Publishing message to NATS Core")

	err = p.conn.Publish(subject, b)
	if err != nil {
		p.log.Error().
			Err(err).
			Str("subject", subject).
			Msg("Failed to publish message to NATS Core")
		return fmt.Errorf("failed to publish to NATS Core: %w", err)
	}

	return nil
}

// PublishEvent publishes a structured event with metadata to JetStream
func (p *Publisher) PublishEvent(subject, eventType string, data interface{}) error {
	event := Event{
		Type:    eventType,
		Version: 1,
		Data:    data,
		Meta: EventMeta{
			Timestamp: time.Now().UnixMilli(),
			Service:   p.serviceName,
		},
	}

	return p.PublishJetStream(subject, event)
}

// PublishEventCore publishes a structured event with metadata using NATS Core
func (p *Publisher) PublishEventCore(subject, eventType string, data interface{}) error {
	event := Event{
		Type:    eventType,
		Version: 1,
		Data:    data,
		Meta: EventMeta{
			Timestamp: time.Now().UnixMilli(),
			Service:   p.serviceName,
		},
	}

	return p.PublishCore(subject, event)
}

// Publish is a convenience method that uses JetStream by default
// For backward compatibility with existing code
func (p *Publisher) Publish(subject string, data interface{}) error {
	return p.PublishJetStream(subject, data)
}

// Close closes the underlying NATS connection
func (p *Publisher) Close() {
	if p.conn != nil {
		p.log.Debug().Msg("Closing NATS Publisher")
		p.conn.Close()
	}
}
