package apikey

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

// Validator valida API keys para comunicación interna
type Validator struct {
	keys map[string]string // key -> service name
}

// Config configuración del sistema de API keys
type Config struct {
	Keys        map[string]string `json:"keys"`         // key -> service name
	HeaderName  string            `json:"header_name"`  // Default: "X-API-Key"
	AllowBearer bool              `json:"allow_bearer"` // Allow Authorization: Bearer
	AllowQuery  bool              `json:"allow_query"`  // Allow ?api_key=
	QueryParam  string            `json:"query_param"`  // Default: "api_key"
}

// DefaultConfig configuración por defecto
func DefaultConfig() *Config {
	return &Config{
		Keys:        make(map[string]string),
		HeaderName:  "X-API-Key",
		AllowBearer: true,
		AllowQuery:  true,
		QueryParam:  "api_key",
	}
}

// NewValidator crea un nuevo validador de API keys
func NewValidator(keys map[string]string) *Validator {
	return &Validator{
		keys: keys,
	}
}

// NewValidatorFromConfig crea un validador desde configuración
func NewValidatorFromConfig(config *Config) *Validator {
	return &Validator{
		keys: config.Keys,
	}
}

// ValidateKey valida una API key y retorna el nombre del servicio
func (v *Validator) ValidateKey(key string) (string, bool) {
	if key == "" {
		return "", false
	}

	serviceName, exists := v.keys[key]

	if !exists {
		// Log solo cuando falla la validación
		maskedKey := key
		if len(key) > 8 {
			maskedKey = key[:4] + "..." + key[len(key)-4:]
		} else {
			maskedKey = "***"
		}

		log.Error().
			Str("api_key_received", maskedKey).
			Int("registered_keys_count", len(v.keys)).
			Msg("API key not found in validator")
	}

	return serviceName, exists
}

// ExtractAPIKey extrae la API key de una request HTTP
func (v *Validator) ExtractAPIKey(r *http.Request, config *Config) string {
	if config == nil {
		config = DefaultConfig()
	}

	// 1. Header personalizado (X-API-Key por defecto)
	if key := r.Header.Get(config.HeaderName); key != "" {
		return key
	}

	// 2. Authorization Bearer (si está habilitado)
	if config.AllowBearer {
		auth := r.Header.Get("Authorization")
		if strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			return strings.TrimSpace(auth[7:])
		}
	}

	// 3. Query parameter (si está habilitado)
	if config.AllowQuery {
		if key := r.URL.Query().Get(config.QueryParam); key != "" {
			return key
		}
	}

	return ""
}

// AddKey agrega una nueva API key al validador
func (v *Validator) AddKey(key, serviceName string) {
	v.keys[key] = serviceName
}

// RemoveKey elimina una API key del validador
func (v *Validator) RemoveKey(key string) {
	delete(v.keys, key)
}

// ListServices retorna una lista de servicios registrados
func (v *Validator) ListServices() []string {
	services := make(map[string]struct{})
	for _, service := range v.keys {
		services[service] = struct{}{}
	}

	result := make([]string, 0, len(services))
	for service := range services {
		result = append(result, service)
	}
	return result
}

// HasService verifica si un servicio tiene al menos una API key registrada
func (v *Validator) HasService(serviceName string) bool {
	for _, service := range v.keys {
		if service == serviceName {
			return true
		}
	}
	return false
}

// GenerateKey genera una nueva API key para un servicio (simple implementation)
func GenerateKey(serviceName string) string {
	// En producción usar crypto/rand para generar keys seguras
	return fmt.Sprintf("%s-key-%d", serviceName, len(serviceName)*12345)
}
