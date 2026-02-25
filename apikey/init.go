package apikey

import (
	"fmt"
	"os"
	"strings"
)

// InitOptions opciones para inicializar el sistema de API keys
type InitOptions struct {
	// RequireAll si true, falla si falta alguna API key estándar
	RequireAll bool

	// Silent si true, no imprime información de estado
	Silent bool

	// AutoGenerate si true, genera claves de desarrollo si faltan
	AutoGenerate bool

	// CheckOnly si true, solo verifica sin crear validador
	CheckOnly bool
}

// DefaultInitOptions configuración por defecto para inicialización
func DefaultInitOptions() *InitOptions {
	return &InitOptions{
		RequireAll:   false,
		Silent:       false,
		AutoGenerate: false,
		CheckOnly:    false,
	}
}

// InitResult resultado de la inicialización del sistema
type InitResult struct {
	Validator     *Validator
	LoadedKeys    []string
	MissingKeys   []string
	GeneratedKeys []string
	Success       bool
	Error         error
}

// InitConnectAPIKeys inicializa el sistema de API keys con opciones flexibles
func InitConnectAPIKeys(options *InitOptions) *InitResult {
	if options == nil {
		options = DefaultInitOptions()
	}

	result := &InitResult{
		LoadedKeys:    []string{},
		MissingKeys:   []string{},
		GeneratedKeys: []string{},
		Success:       false,
	}

	// Cargar claves existentes
	loadedKeys := loadExistingKeys(result)

	// Generar claves faltantes si es necesario
	if options.AutoGenerate {
		generateMissingKeys(result, loadedKeys)
	}

	// Validar configuración
	if err := validateConfiguration(options, result); err != nil {
		result.Error = err
		return result
	}

	// Finalizar inicialización
	finalizeInit(options, result, loadedKeys)

	return result
}

// loadExistingKeys carga las API keys existentes desde variables de entorno
func loadExistingKeys(result *InitResult) map[string]string {
	loadedKeys := make(map[string]string)

	for service, envVar := range DefaultEnvMapping {
		apiKey := os.Getenv(envVar)
		if apiKey != "" {
			loadedKeys[apiKey] = service
			result.LoadedKeys = append(result.LoadedKeys, fmt.Sprintf("%s=%s", envVar, maskAPIKey(apiKey)))
		} else {
			result.MissingKeys = append(result.MissingKeys, envVar)
		}
	}

	return loadedKeys
}

// generateMissingKeys genera claves de desarrollo para las que faltan
func generateMissingKeys(result *InitResult, loadedKeys map[string]string) {
	if len(result.MissingKeys) == 0 {
		return
	}

	var remainingMissing []string

	for _, envVar := range result.MissingKeys {
		service := getServiceFromEnvVar(envVar)
		if service != "" {
			generatedKey := fmt.Sprintf("dev-%s-key-%d", service, len(service)*12345)
			if err := os.Setenv(envVar, generatedKey); err == nil {
				loadedKeys[generatedKey] = "connect-" + service
				result.GeneratedKeys = append(result.GeneratedKeys, envVar)
			} else {
				remainingMissing = append(remainingMissing, envVar)
			}
		} else {
			remainingMissing = append(remainingMissing, envVar)
		}
	}

	result.MissingKeys = remainingMissing
}

// validateConfiguration valida que la configuración sea válida
func validateConfiguration(options *InitOptions, result *InitResult) error {
	if options.RequireAll && len(result.MissingKeys) > 0 {
		return fmt.Errorf("missing required API keys: %v", result.MissingKeys)
	}
	return nil
}

// finalizeInit completa la inicialización
func finalizeInit(options *InitOptions, result *InitResult, loadedKeys map[string]string) {
	if !options.CheckOnly {
		result.Validator = NewValidator(loadedKeys)
	}

	result.Success = len(loadedKeys) > 0 || options.CheckOnly

	if !options.Silent {
		printInitResult(result)
	}
}

// QuickInit inicialización rápida para desarrollo
func QuickInit() (*Validator, error) {
	result := InitConnectAPIKeys(&InitOptions{
		RequireAll:   false,
		Silent:       false,
		AutoGenerate: true,
		CheckOnly:    false,
	})

	return result.Validator, result.Error
}

// ProductionInit inicialización estricta para producción
func ProductionInit() (*Validator, error) {
	result := InitConnectAPIKeys(&InitOptions{
		RequireAll:   true,
		Silent:       true,
		AutoGenerate: false,
		CheckOnly:    false,
	})

	return result.Validator, result.Error
}

// CheckSetup verifica la configuración sin crear validador
func CheckSetup() error {
	result := InitConnectAPIKeys(&InitOptions{
		RequireAll: true,
		Silent:     false,
		CheckOnly:  true,
	})

	return result.Error
}

// Helper functions

func getServiceFromEnvVar(envVar string) string {
	for service, env := range DefaultEnvMapping {
		if env == envVar {
			return strings.TrimPrefix(service, "connect-")
		}
	}
	return ""
}

func printInitResult(result *InitResult) {
	fmt.Println("🔐 Connect API Keys Initialization")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	if len(result.LoadedKeys) > 0 {
		fmt.Printf("✅ Loaded Keys (%d):\n", len(result.LoadedKeys))
		for _, key := range result.LoadedKeys {
			fmt.Printf("   %s\n", key)
		}
	}

	if len(result.GeneratedKeys) > 0 {
		fmt.Printf("🔧 Generated Keys (%d):\n", len(result.GeneratedKeys))
		for _, key := range result.GeneratedKeys {
			fmt.Printf("   %s (dev mode)\n", key)
		}
	}

	if len(result.MissingKeys) > 0 {
		fmt.Printf("⚠️ Missing Keys (%d):\n", len(result.MissingKeys))
		for _, key := range result.MissingKeys {
			fmt.Printf("   %s\n", key)
		}
		fmt.Println("\n💡 Add these to your .env file or use AutoGenerate for dev mode")
	}

	if result.Success {
		fmt.Println("✅ API Key system ready")
	} else {
		fmt.Println("❌ API Key system initialization failed")
		if result.Error != nil {
			fmt.Printf("   Error: %v\n", result.Error)
		}
	}
	fmt.Println()
}
