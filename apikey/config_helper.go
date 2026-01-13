package apikey

import (
	"os"
	"path/filepath"
)

// ConfigHelper helpers opcionales para cargar configuración en backends
type ConfigHelper struct {
	EnvFiles []string
	Required []string
}

// DefaultConfigHelper configuración por defecto para servicios Connect
func DefaultConfigHelper() *ConfigHelper {
	return &ConfigHelper{
		EnvFiles: []string{".env.development", ".env.local", ".env"},
		Required: []string{"AUTH_API_KEY", "CORE_API_KEY", "LOBBY_API_KEY", "RT_API_KEY"},
	}
}

// LoadEnvFiles intenta cargar archivos .env (helper opcional para backends)
// Nota: Requiere que el backend incluya github.com/joho/godotenv
func (h *ConfigHelper) LoadEnvFiles() error {
	// Nota: Este es solo un helper, el backend sigue siendo responsable de la carga
	for _, envFile := range h.EnvFiles {
		if _, err := os.Stat(envFile); err == nil {
			// Backend necesita llamar godotenv.Load(envFile) directamente
			// Este helper solo ayuda a decidir qué archivo cargar
			return nil
		}
	}

	return nil
}

// FindEnvFile encuentra el primer archivo .env disponible
func (h *ConfigHelper) FindEnvFile() string {
	for _, envFile := range h.EnvFiles {
		if absPath, err := filepath.Abs(envFile); err == nil {
			if _, err := os.Stat(absPath); err == nil {
				return absPath
			}
		}
	}
	return ""
}

// ValidateRequired verifica que las variables requeridas estén presentes
func (h *ConfigHelper) ValidateRequired() []string {
	var missing []string

	for _, envVar := range h.Required {
		if os.Getenv(envVar) == "" {
			missing = append(missing, envVar)
		}
	}

	return missing
}

// Usage example for backend services
/*
Example usage in Connect-Core/cmd/server/main.go:

import (
    "github.com/joho/godotenv"
    apikey "github.com/AoC-Gamers/connect-libraries/apikey"
)

func main() {
    // 1. Backend loads configuration
    helper := apikey.DefaultConfigHelper()

    if envFile := helper.FindEnvFile(); envFile != "" {
        if err := godotenv.Load(envFile); err != nil {
            log.Printf("Failed to load %s: %v", envFile, err)
        }
    }

    // 2. Validate required variables
    if missing := helper.ValidateRequired(); len(missing) > 0 {
        log.Fatalf("Missing required environment variables: %v", missing)
    }

    // 3. Package reads from loaded environment
    validator, err := apikey.LoadConnectAPIKeys()
    if err != nil {
        log.Fatal("API key setup failed:", err)
    }

    // 4. Use in application
    router.Use(ginapi.RequireAPIKey(validator))
}
*/
