package apikey

import (
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
)

// ServiceAPIKeys mapas estándar de servicios Connect
var (
	// ConnectServices lista de servicios conocidos
	ConnectServices = []string{
		"connect-auth",
		"connect-core",
		"connect-lobby",
		"connect-rt",
	}

	// DefaultEnvMapping mapeo de servicio -> variable de entorno
	DefaultEnvMapping = map[string]string{
		"connect-auth":  "AUTH_API_KEY",
		"connect-core":  "CORE_API_KEY",
		"connect-lobby": "LOBBY_API_KEY",
		"connect-rt":    "RT_API_KEY",
	}
)

// EnvConfig configuración basada en variables de entorno
type EnvConfig struct {
	// ServiceMapping mapea nombres de servicio a variables de entorno
	ServiceMapping map[string]string

	// AllowMissing si true, no falla si falta alguna API key
	AllowMissing bool

	// Prefix prefijo opcional para las variables de entorno
	Prefix string

	// CustomKeys API keys adicionales no estándar
	CustomKeys map[string]string
}

// DefaultEnvConfig configuración por defecto para servicios Connect
func DefaultEnvConfig() *EnvConfig {
	return &EnvConfig{
		ServiceMapping: DefaultEnvMapping,
		AllowMissing:   false,
		Prefix:         "",
		CustomKeys:     make(map[string]string),
	}
}

// NewValidatorFromEnv crea un validador cargando API keys desde variables de entorno
func NewValidatorFromEnv(config *EnvConfig) (*Validator, error) {
	if config == nil {
		config = DefaultEnvConfig()
	}

	keys := make(map[string]string)
	var missingKeys []string

	// Cargar API keys de servicios estándar
	for service, envVar := range config.ServiceMapping {
		fullEnvVar := config.Prefix + envVar
		apiKey := os.Getenv(fullEnvVar)

		if apiKey == "" {
			if !config.AllowMissing {
				missingKeys = append(missingKeys, fullEnvVar)
			}
			continue
		}

		keys[apiKey] = service
	}

	// Cargar API keys personalizadas
	for envVar, service := range config.CustomKeys {
		fullEnvVar := config.Prefix + envVar
		apiKey := os.Getenv(fullEnvVar)

		if apiKey != "" {
			keys[apiKey] = service
		} else if !config.AllowMissing {
			missingKeys = append(missingKeys, fullEnvVar)
		}
	}

	// Reportar claves faltantes si es necesario
	if len(missingKeys) > 0 && !config.AllowMissing {
		log.Error().
			Strs("missing_keys", missingKeys).
			Msg("❌ Missing required API key environment variables")
		return nil, fmt.Errorf("missing required API key environment variables: %v", missingKeys)
	}

	return NewValidator(keys), nil
}

// LoadConnectAPIKeys carga las API keys estándar de servicios Connect desde env
func LoadConnectAPIKeys() (*Validator, error) {
	return NewValidatorFromEnv(DefaultEnvConfig())
}

// LoadConnectAPIKeysPermissive carga API keys permitiendo que falten algunas
func LoadConnectAPIKeysPermissive() (*Validator, error) {
	config := DefaultEnvConfig()
	config.AllowMissing = true
	return NewValidatorFromEnv(config)
}

// GetServiceAPIKey obtiene la API key de un servicio específico desde env
func GetServiceAPIKey(service string) string {
	if envVar, exists := DefaultEnvMapping[service]; exists {
		return os.Getenv(envVar)
	}
	return ""
}

// ValidateEnvSetup verifica que todas las variables de entorno necesarias estén configuradas
func ValidateEnvSetup() error {
	var missing []string

	for service, envVar := range DefaultEnvMapping {
		if os.Getenv(envVar) == "" {
			missing = append(missing, fmt.Sprintf("%s (for %s)", envVar, service))
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required API key environment variables:\n%v\n\nPlease ensure your .env file contains these keys", missing)
	}

	return nil
}

// PrintEnvStatus muestra el estado de las API keys configuradas
func PrintEnvStatus() {
	fmt.Println("🔐 Connect API Keys Status:")

	for service, envVar := range DefaultEnvMapping {
		apiKey := os.Getenv(envVar)
		if apiKey == "" {
			fmt.Printf("  ❌ %s: %s (not set)\n", service, envVar)
		} else {
			// Mostrar solo los primeros y últimos 4 caracteres por seguridad
			masked := maskAPIKey(apiKey)
			fmt.Printf("  ✅ %s: %s (%s)\n", service, envVar, masked)
		}
	}
}

// maskAPIKey enmascara una API key para logging seguro
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + "****" + key[len(key)-4:]
}

// GenerateDevAPIKeys genera API keys para desarrollo basadas en patrones estándar
func GenerateDevAPIKeys() map[string]string {
	keys := make(map[string]string)

	for service := range DefaultEnvMapping {
		keys[service] = fmt.Sprintf("%s-dev-%s", service, "internal-key-change-in-production")
	}

	return keys
}

// WithCustomService agrega un servicio personalizado a la configuración
func (c *EnvConfig) WithCustomService(serviceName, envVar string) *EnvConfig {
	if c.CustomKeys == nil {
		c.CustomKeys = make(map[string]string)
	}
	c.CustomKeys[envVar] = serviceName
	return c
}

// WithPrefix establece un prefijo para todas las variables de entorno
func (c *EnvConfig) WithPrefix(prefix string) *EnvConfig {
	c.Prefix = prefix
	return c
}

// AllowMissingKeys permite que falten algunas API keys sin generar error
func (c *EnvConfig) AllowMissingKeys() *EnvConfig {
	c.AllowMissing = true
	return c
}
