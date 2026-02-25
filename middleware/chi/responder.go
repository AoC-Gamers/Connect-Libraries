package chi

import (
	"encoding/json"
	"net/http"
)

// ErrorCode representa códigos de error estandarizados
type ErrorCode string

const (
	// CodeUnauthorized indica que la petición requiere autenticación
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"

	// CodeTokenExpired indica que el JWT ha expirado
	CodeTokenExpired ErrorCode = "TOKEN_EXPIRED"

	// CodePolicyVersionMismatch indica mismatch entre token y política actual
	CodePolicyVersionMismatch ErrorCode = "POLICY_VERSION_MISMATCH"

	// CodeInsufficientPermissions indica permisos insuficientes en general
	CodeInsufficientPermissions ErrorCode = "INSUFFICIENT_PERMISSIONS"
)

// ErrorResponse representa una respuesta de error estructurada
// Basado en RFC 7807 (Problem Details for HTTP APIs) simplificado
type ErrorResponse struct {
	// Error es el mensaje corto y legible
	Error string `json:"error"`

	// Code es el código de error programático para el frontend
	Code ErrorCode `json:"code,omitempty"`

	// Status es el código HTTP (redundante pero útil para debugging)
	Status int `json:"status"`

	// Detail es una explicación detallada del error
	Detail string `json:"detail,omitempty"`

	// Meta contiene metadata adicional específica del contexto
	Meta map[string]interface{} `json:"meta,omitempty"`
}

// ErrorResponder permite desacoplar las respuestas de error
type ErrorResponder interface {
	Unauthorized(w http.ResponseWriter, detail string)
	TokenExpired(w http.ResponseWriter)
	PolicyVersionMismatch(w http.ResponseWriter, tokenVersion, currentVersion int)
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

// TokenExpired responde cuando el JWT ha expirado
func (DefaultErrorResponder) TokenExpired(w http.ResponseWriter) {
	respondError(w, http.StatusUnauthorized, CodeTokenExpired,
		"unauthorized",
		"JWT token has expired",
		map[string]interface{}{
			"should_refresh": true,
		})
}

// PolicyVersionMismatch responde cuando hay mismatch de policy version
func (DefaultErrorResponder) PolicyVersionMismatch(w http.ResponseWriter, tokenVersion, currentVersion int) {
	respondError(w, http.StatusUnauthorized, CodePolicyVersionMismatch,
		"policy version mismatch",
		"Token policy version does not match current version",
		map[string]interface{}{
			"token_version":   tokenVersion,
			"current_version": currentVersion,
			"should_reauth":   true,
		})
}

// InsufficientPermissions responde con permisos insuficientes genérico
func (DefaultErrorResponder) InsufficientPermissions(w http.ResponseWriter, action string) {
	respondWithDetail(w, http.StatusForbidden, CodeInsufficientPermissions,
		"insufficient permissions",
		"You don't have permission to "+action)
}

const (
	responderHeaderContentType = "Content-Type"
	responderMimeJSON          = "application/json"
)

func respondError(w http.ResponseWriter, status int, code ErrorCode, message, detail string, meta map[string]interface{}) {
	resp := ErrorResponse{
		Error:  message,
		Code:   code,
		Status: status,
		Detail: detail,
		Meta:   meta,
	}

	w.Header().Set(responderHeaderContentType, responderMimeJSON)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func respondWithDetail(w http.ResponseWriter, status int, code ErrorCode, message, detail string) {
	respondError(w, status, code, message, detail, nil)
}
