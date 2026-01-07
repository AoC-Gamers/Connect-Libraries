package errors

import (
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

// InternalErrorCode códigos estándar para servicios internos
type InternalErrorCode string

const (
	// Authentication & Authorization
	InternalUnauthorized InternalErrorCode = "INTERNAL_UNAUTHORIZED"
	InternalForbidden    InternalErrorCode = "INTERNAL_FORBIDDEN"

	// Resource Management
	InternalNotFound InternalErrorCode = "INTERNAL_NOT_FOUND"
	InternalConflict InternalErrorCode = "INTERNAL_CONFLICT"

	// Validation & Input
	InternalValidation InternalErrorCode = "INTERNAL_VALIDATION"
	InternalBadRequest InternalErrorCode = "INTERNAL_BAD_REQUEST"

	// System & Infrastructure
	InternalDatabase    InternalErrorCode = "INTERNAL_DATABASE"
	InternalTimeout     InternalErrorCode = "INTERNAL_TIMEOUT"
	InternalServiceDown InternalErrorCode = "INTERNAL_SERVICE_DOWN"
	InternalRateLimit   InternalErrorCode = "INTERNAL_RATE_LIMIT"
	InternalAuthzCheck  InternalErrorCode = "INTERNAL_AUTHZ_CHECK"

	// Generic
	InternalServerError InternalErrorCode = "INTERNAL_SERVER_ERROR"
)

// InternalErrorResponse estructura estándar para respuestas entre servicios
type InternalErrorResponse struct {
	Code    InternalErrorCode `json:"code"`
	Message string            `json:"message"`
	Service string            `json:"service"`
	Details interface{}       `json:"details,omitempty"`
	Status  int               `json:"status"`
}

// detectServiceName intenta detectar el nombre del servicio desde un request
func detectServiceName(r *http.Request) string {
	// Intentar obtener desde header personalizado
	if service := r.Header.Get("X-Service-Name"); service != "" {
		return service
	}

	// Detectar por path
	path := r.URL.Path
	if len(path) > 0 {
		if path[0] == '/' {
			path = path[1:]
		}
		if len(path) > 0 {
			parts := strings.Split(path, "/")
			if len(parts) > 0 {
				switch parts[0] {
				case "auth":
					return "connect-auth"
				case "core":
					return "connect-core"
				case "lobby":
					return "connect-lobby"
				case "rt":
					return "connect-rt"
				}
			}
		}
	}

	// Fallback genérico
	return "connect-service"
}

// RespondInternalServiceError respuesta estándar para errores internos entre servicios
func RespondInternalServiceError(w http.ResponseWriter, r *http.Request, code InternalErrorCode, message string, details interface{}) {
	service := detectServiceName(r)

	statusCode := getInternalStatusCode(code)

	response := InternalErrorResponse{
		Code:    code,
		Message: message,
		Service: service,
		Details: details,
		Status:  statusCode,
	}

	// Log estructurado para debugging interno
	log.Error().
		Str("internal_code", string(code)).
		Str("service", service).
		Str("path", r.URL.Path).
		Str("method", r.Method).
		Int("status", statusCode).
		Interface("details", details).
		Msg("Internal service error")

	render.Status(r, statusCode)
	render.JSON(w, r, response)
}

// ============================================
// HELPERS ESPECÍFICOS PARA ERRORES INTERNOS
// ============================================

// RespondInternalUnauthorized (401) - API key inválida o faltante
func RespondInternalUnauthorized(w http.ResponseWriter, r *http.Request, message string) {
	if message == "" {
		message = "invalid or missing API key"
	}
	RespondInternalServiceError(w, r, InternalUnauthorized, message, nil)
}

// RespondInternalForbidden (403) - API key válida pero servicio no autorizado
func RespondInternalForbidden(w http.ResponseWriter, r *http.Request, allowedServices []string, actualService string) {
	details := map[string]interface{}{
		"allowed_services": allowedServices,
		"actual_service":   actualService,
	}
	RespondInternalServiceError(w, r, InternalForbidden, "service not authorized for this endpoint", details)
}

// RespondInternalNotFound (404) - Recurso no encontrado
func RespondInternalNotFound(w http.ResponseWriter, r *http.Request, resource string, id string) {
	details := map[string]interface{}{
		"resource": resource,
		"id":       id,
	}
	RespondInternalServiceError(w, r, InternalNotFound, "resource not found", details)
}

// RespondInternalValidation (400) - Error de validación
func RespondInternalValidation(w http.ResponseWriter, r *http.Request, field string, reason string) {
	details := map[string]interface{}{
		"field":  field,
		"reason": reason,
	}
	RespondInternalServiceError(w, r, InternalValidation, "validation error", details)
}

// RespondInternalBadRequest (400) - Request malformado
func RespondInternalBadRequest(w http.ResponseWriter, r *http.Request, reason string) {
	details := map[string]interface{}{
		"reason": reason,
	}
	RespondInternalServiceError(w, r, InternalBadRequest, "bad request", details)
}

// RespondInternalDatabase (500) - Error de base de datos
func RespondInternalDatabase(w http.ResponseWriter, r *http.Request, operation string, err error) {
	details := map[string]interface{}{
		"operation": operation,
		"error":     err.Error(),
	}
	RespondInternalServiceError(w, r, InternalDatabase, "database operation failed", details)
}

// RespondInternalTimeout (500) - Timeout en operación
func RespondInternalTimeout(w http.ResponseWriter, r *http.Request, operation string, timeout string) {
	details := map[string]interface{}{
		"operation": operation,
		"timeout":   timeout,
	}
	RespondInternalServiceError(w, r, InternalTimeout, "operation timed out", details)
}

// RespondInternalServiceDown (503) - Servicio dependiente no disponible
func RespondInternalServiceDown(w http.ResponseWriter, r *http.Request, service string) {
	details := map[string]interface{}{
		"unavailable_service": service,
	}
	RespondInternalServiceError(w, r, InternalServiceDown, "dependent service unavailable", details)
}

// ============================================
// UTILITY FUNCTIONS
// ============================================

// getInternalStatusCode mapea códigos internos a códigos HTTP
func getInternalStatusCode(code InternalErrorCode) int {
	switch code {
	case InternalUnauthorized:
		return http.StatusUnauthorized
	case InternalForbidden:
		return http.StatusForbidden
	case InternalNotFound:
		return http.StatusNotFound
	case InternalConflict:
		return http.StatusConflict
	case InternalValidation, InternalBadRequest:
		return http.StatusBadRequest
	case InternalRateLimit:
		return http.StatusTooManyRequests
	case InternalServiceDown:
		return http.StatusServiceUnavailable
	case InternalDatabase, InternalTimeout, InternalServerError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}
