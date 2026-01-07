package errors

import (
	"encoding/json"
	"net/http"
)

// ErrorResponse representa una respuesta de error estructurada
// Basado en RFC 7807 (Problem Details for HTTP APIs) simplificado
type ErrorResponse struct {
	// Error es el mensaje corto y legible (mantiene compatibilidad con APIs existentes)
	Error string `json:"error"`

	// Code es el código de error programático para el frontend
	Code ErrorCode `json:"code,omitempty"`

	// Status es el código HTTP (redundante pero útil para debugging)
	Status int `json:"status"`

	// Detail es una explicación detallada del error
	Detail string `json:"detail,omitempty"`

	// Meta contiene metadata adicional específica del contexto
	// Por ejemplo: scope_id, required_permission, field_name, etc.
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// RespondError escribe una respuesta de error estructurada
// Es la función principal para enviar errores al cliente
func RespondError(w http.ResponseWriter, status int, code ErrorCode, message, detail string, meta map[string]interface{}) {
	resp := ErrorResponse{
		Error:  message,
		Code:   code,
		Status: status,
		Detail: detail,
		Meta:   meta,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

// RespondSimpleError es un helper para errores sin metadata
// Útil para migraciones graduales desde APIs existentes
func RespondSimpleError(w http.ResponseWriter, status int, code ErrorCode, message string) {
	RespondError(w, status, code, message, "", nil)
}

// RespondWithDetail es un helper para errores con detalle pero sin metadata
func RespondWithDetail(w http.ResponseWriter, status int, code ErrorCode, message, detail string) {
	RespondError(w, status, code, message, detail, nil)
}

// RespondLegacyError mantiene compatibilidad con el formato antiguo {"error": "message"}
// Útil durante la transición para no romper clientes existentes
func RespondLegacyError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// RespondJSON es un helper para respuestas exitosas
// Incluido aquí para mantener consistencia en el paquete
func RespondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		_ = json.NewEncoder(w).Encode(data)
	}
}
