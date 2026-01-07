package connectnats

import (
	"fmt"
	"time"

	natsio "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// StreamConfig define la configuración de un stream de JetStream
type StreamConfig struct {
	Name        string
	Subjects    []string
	Description string
	MaxAge      time.Duration
	MaxBytes    int64
	Replicas    int
}

// DefaultStreamConfigs retorna las configuraciones de streams por defecto para Connect
func DefaultStreamConfigs() []StreamConfig {
	return []StreamConfig{
		{
			Name:        "CONNECT_EVENTS",
			Subjects:    []string{"connect.>"},
			Description: "Main stream for all Connect system events (invitations, notifications, presence)",
			MaxAge:      7 * 24 * time.Hour, // 7 días
			MaxBytes:    1024 * 1024 * 1024, // 1 GB
			Replicas:    1,
		},
		{
			Name:        "CACHE",
			Subjects:    []string{"cache.>"},
			Description: "Stream for cache invalidation events (ephemeral)",
			MaxAge:      5 * time.Minute, // TTL corto - eventos efímeros
			MaxBytes:    1 * 1024 * 1024, // 1MB máximo
			Replicas:    1,
		},
	}
}

// EnsureStreams crea o actualiza streams de JetStream si no existen
func EnsureStreams(conn *natsio.Conn, configs []StreamConfig) error {
	if conn == nil {
		return fmt.Errorf("nil NATS connection")
	}

	js, err := conn.JetStream()
	if err != nil {
		return fmt.Errorf("failed to get JetStream context: %w", err)
	}

	for _, cfg := range configs {
		if err := ensureStream(js, cfg); err != nil {
			log.Warn().
				Err(err).
				Str("stream", cfg.Name).
				Msg("Failed to ensure JetStream stream - continuing anyway")
		}
	}

	return nil
}

// ensureStream crea o actualiza un stream individual
func ensureStream(js natsio.JetStreamContext, cfg StreamConfig) error {
	streamCfg := &natsio.StreamConfig{
		Name:        cfg.Name,
		Subjects:    cfg.Subjects,
		Description: cfg.Description,
		MaxAge:      cfg.MaxAge,
		MaxBytes:    cfg.MaxBytes,
		Storage:     natsio.FileStorage,
		Replicas:    cfg.Replicas,
		Retention:   natsio.LimitsPolicy,
		Discard:     natsio.DiscardOld,
	}

	// Intentar obtener info del stream
	info, err := js.StreamInfo(cfg.Name)
	if err != nil {
		// Stream no existe, crearlo
		if err == natsio.ErrStreamNotFound {
			_, err := js.AddStream(streamCfg)
			if err != nil {
				return fmt.Errorf("failed to create stream %s: %w", cfg.Name, err)
			}
			return nil
		}
		return fmt.Errorf("failed to check stream %s: %w", cfg.Name, err)
	}

	// Stream exists (silent check)
	_ = info
	return nil
}

// ConnectWithJetStream establece conexión con NATS y configura JetStream automáticamente
func ConnectWithJetStream(cfg *Config, streamConfigs []StreamConfig) (*natsio.Conn, error) {
	// Conectar primero
	conn, err := Connect(cfg)
	if err != nil {
		return nil, err
	}

	// Configurar streams si se proporcionaron
	if len(streamConfigs) == 0 {
		streamConfigs = DefaultStreamConfigs()
	}

	if err := EnsureStreams(conn, streamConfigs); err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to ensure JetStream streams - connection will continue")
	}

	return conn, nil
}

// MustConnectWithJetStream conecta a NATS con JetStream y retorna nil si falla
func MustConnectWithJetStream(cfg *Config, streamConfigs []StreamConfig) *natsio.Conn {
	conn, err := ConnectWithJetStream(cfg, streamConfigs)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("NATS+JetStream connection failed - service will continue without real-time events")
		return nil
	}
	return conn
}
