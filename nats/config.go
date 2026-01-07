package connectnats

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	natsio "github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

// Config contiene la configuración estándar para NATS en todos los servicios
type Config struct {
	// URL del servidor NATS (ej: "nats://localhost:4222")
	URL string

	// ClientID identifica la aplicación conectada
	ClientID string

	// Tiempo de espera entre reconexiones (segundos)
	ReconnectWait int

	// Máximo número de intentos de reconexión (-1 = infinito)
	MaxReconnects int

	// Timeout para operaciones de conexión
	Timeout time.Duration
}

// DefaultConfig retorna la configuración por defecto para NATS
func DefaultConfig() *Config {
	return &Config{
		URL:           getEnvOrDefault("NATS_URL", "nats://localhost:4222"),
		ClientID:      getEnvOrDefault("NATS_CLIENT_ID", "connect-service"),
		ReconnectWait: 2,
		MaxReconnects: -1, // Infinito
		Timeout:       10 * time.Second,
	}
}

// Connect establece conexión con NATS usando configuración estandarizada
// Soporta autenticación por token, credenciales y TLS
func Connect(cfg *Config) (*natsio.Conn, error) {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	opts := buildConnectionOptions(cfg)

	// Intentar conexión
	natsConn, err := natsio.Connect(cfg.URL, opts...)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to connect to NATS - continuing without real-time events")
		return nil, err
	}

	return natsConn, nil
}

// MustConnect conecta a NATS y retorna nil si falla (no termina el programa)
// Permite que el servicio continúe sin eventos en tiempo real
func MustConnect(cfg *Config) *natsio.Conn {
	conn, err := Connect(cfg)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("NATS connection failed - service will continue without real-time events")
		return nil
	}
	return conn
}

// buildConnectionOptions construye las opciones de conexión estándar
func buildConnectionOptions(cfg *Config) []natsio.Option {
	opts := []natsio.Option{
		natsio.Name(cfg.ClientID),
		natsio.Timeout(cfg.Timeout),
		natsio.ReconnectWait(time.Duration(cfg.ReconnectWait) * time.Second),
		natsio.MaxReconnects(cfg.MaxReconnects),
		natsio.DisconnectErrHandler(func(nc *natsio.Conn, err error) {
			if err != nil {
				log.Warn().
					Err(err).
					Str("client_id", cfg.ClientID).
					Msg("NATS disconnected")
			}
		}),
		natsio.ReconnectHandler(func(nc *natsio.Conn) {
			// Reconnected silently
		}),
		natsio.ClosedHandler(func(nc *natsio.Conn) {
			log.Warn().
				Str("client_id", cfg.ClientID).
				Msg("NATS connection closed")
		}),
	}

	// Autenticación por token (más común en desarrollo)
	if token := os.Getenv("NATS_TOKEN"); token != "" {
		opts = append(opts, natsio.Token(token))
	}

	// Autenticación por credenciales (.creds file)
	if credFile := os.Getenv("NATS_CREDS_FILE"); credFile != "" {
		log.Debug().Str("creds_file", credFile).Msg("Using NATS credentials file")
		opts = append(opts, natsio.UserCredentials(credFile))
	}

	// Autenticación por usuario/contraseña
	if user := os.Getenv("NATS_USER"); user != "" {
		password := os.Getenv("NATS_PASSWORD")
		log.Debug().Str("user", user).Msg("Using NATS user/password authentication")
		opts = append(opts, natsio.UserInfo(user, password))
	}

	// TLS opcional
	if cert := os.Getenv("NATS_TLS_CERT"); cert != "" {
		key := os.Getenv("NATS_TLS_KEY")
		ca := os.Getenv("NATS_TLS_CA")
		if tlsCfg, tlsErr := createTLSConfig(cert, key, ca); tlsErr == nil {
			log.Debug().Msg("Using NATS TLS configuration")
			opts = append(opts, natsio.Secure(tlsCfg))
		} else {
			log.Warn().
				Err(tlsErr).
				Msg("Failed to build TLS config for NATS - continuing without TLS")
		}
	}

	return opts
}

// createTLSConfig construye configuración TLS desde archivos de certificados
func createTLSConfig(certFile, keyFile, caFile string) (*tls.Config, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load client cert/key: %w", err)
	}

	caCert, err := os.ReadFile(caFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		return nil, fmt.Errorf("failed to parse CA cert")
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
		MinVersion:   tls.VersionTLS12,
	}, nil
}

// getEnvOrDefault retorna el valor de una variable de entorno o un default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
