package apikey

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorCode representa códigos de error estandarizados
type ErrorCode string

const (
	// CodeUnauthorized indica que la petición requiere autenticación
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"

	// CodeInsufficientPermissions indica permisos insuficientes en general
	CodeInsufficientPermissions ErrorCode = "INSUFFICIENT_PERMISSIONS"
)

// ErrorResponse representa una respuesta de error estructurada
type ErrorResponse struct {
	Error  string    `json:"error"`
	Code   ErrorCode `json:"code,omitempty"`
	Status int       `json:"status"`
	Detail string    `json:"detail,omitempty"`
}

// ErrorResponder permite desacoplar las respuestas de error
type ErrorResponder interface {
	Unauthorized(w http.ResponseWriter, detail string)
	InsufficientPermissions(w http.ResponseWriter, action string)
}

// DefaultErrorResponder implementa respuestas estándar JSON
type DefaultErrorResponder struct{}

// ensureResponder retorna un responder válido
func ensureResponder(responder ErrorResponder) ErrorResponder {
	if responder == nil {
		return DefaultErrorResponder{}
	}
	return responder
}

// Unauthorized responde con 401 genérico
func (DefaultErrorResponder) Unauthorized(w http.ResponseWriter, detail string) {
	respondWithDetail(w, http.StatusUnauthorized, CodeUnauthorized, "unauthorized", detail)
}

// InsufficientPermissions responde con permisos insuficientes
func (DefaultErrorResponder) InsufficientPermissions(w http.ResponseWriter, action string) {
	respondWithDetail(w, http.StatusForbidden, CodeInsufficientPermissions,
		"insufficient permissions",
		"You don't have permission to "+action)
}

const (
	responderHeaderContentType = "Content-Type"
	responderMimeJSON          = "application/json"
)

func respondError(w http.ResponseWriter, status int, code ErrorCode, message, detail string) {
	resp := ErrorResponse{
		Error:  message,
		Code:   code,
		Status: status,
		Detail: detail,
	}

	w.Header().Set(responderHeaderContentType, responderMimeJSON)
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode error response: %v", err)
	}
}

func respondWithDetail(w http.ResponseWriter, status int, code ErrorCode, message, detail string) {
	respondError(w, status, code, message, detail)
}
