package errors

// ErrorCode representa códigos de error estandarizados
// Siguiendo el patrón SCREAMING_SNAKE_CASE para consistencia con REST APIs
type ErrorCode string

const (
	// ============================================
	// Authentication & Authorization Errors
	// ============================================

	// CodeUnauthorized indica que la petición requiere autenticación
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"

	// CodeTokenExpired indica que el JWT ha expirado
	CodeTokenExpired ErrorCode = "TOKEN_EXPIRED"

	// CodeTokenInvalid indica que el JWT es inválido o corrupto
	CodeTokenInvalid ErrorCode = "TOKEN_INVALID"

	// CodePolicyVersionMismatch indica mismatch entre token y política actual
	CodePolicyVersionMismatch ErrorCode = "POLICY_VERSION_MISMATCH"

	// CodePermissionDenied indica falta de permisos específicos
	CodePermissionDenied ErrorCode = "PERMISSION_DENIED"

	// CodeInsufficientPermissions indica permisos insuficientes en general
	CodeInsufficientPermissions ErrorCode = "INSUFFICIENT_PERMISSIONS"

	// ============================================
	// Resource Errors
	// ============================================

	// CodeNotFound indica que el recurso solicitado no existe
	CodeNotFound ErrorCode = "NOT_FOUND"

	// CodeMembershipNotFound indica que el usuario no tiene membership en el scope
	CodeMembershipNotFound ErrorCode = "MEMBERSHIP_NOT_FOUND"

	// CodeAlreadyExists indica que el recurso ya existe (conflicto de creación)
	CodeAlreadyExists ErrorCode = "ALREADY_EXISTS"

	// CodeConflict indica conflicto de estado del recurso
	CodeConflict ErrorCode = "CONFLICT"

	// CodeGone indica que el recurso ya no está disponible permanentemente (HTTP 410)
	CodeGone ErrorCode = "GONE"

	// CodeResourceLocked indica que el recurso está bloqueado
	CodeResourceLocked ErrorCode = "RESOURCE_LOCKED"

	// ============================================
	// Validation Errors
	// ============================================

	// CodeValidationError indica error de validación de datos
	CodeValidationError ErrorCode = "VALIDATION_ERROR"

	// CodeInvalidRequest indica petición malformada
	CodeInvalidRequest ErrorCode = "INVALID_REQUEST"

	// CodeBadRequest indica petición incorrecta en general
	CodeBadRequest ErrorCode = "BAD_REQUEST"

	// CodeMissingRequiredField indica campo requerido faltante
	CodeMissingRequiredField ErrorCode = "MISSING_REQUIRED_FIELD"

	// CodeInvalidFormat indica formato de datos incorrecto
	CodeInvalidFormat ErrorCode = "INVALID_FORMAT"

	// CodeOutOfRange indica valor fuera del rango permitido
	CodeOutOfRange ErrorCode = "OUT_OF_RANGE"

	// ============================================
	// Server Errors
	// ============================================

	// CodeInternalError indica error interno del servidor
	CodeInternalError ErrorCode = "INTERNAL_ERROR"

	// CodeDatabaseError indica error de base de datos
	CodeDatabaseError ErrorCode = "DATABASE_ERROR"

	// CodeServiceUnavailable indica servicio temporalmente no disponible
	CodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"

	// CodeTimeout indica que la operación excedió el tiempo límite
	CodeTimeout ErrorCode = "TIMEOUT"

	// ============================================
	// Business Logic Errors
	// ============================================

	// CodeOperationNotAllowed indica operación no permitida por reglas de negocio
	CodeOperationNotAllowed ErrorCode = "OPERATION_NOT_ALLOWED"

	// CodeQuotaExceeded indica que se excedió una cuota o límite
	CodeQuotaExceeded ErrorCode = "QUOTA_EXCEEDED"

	// CodeRateLimitExceeded indica exceso de rate limit
	CodeRateLimitExceeded ErrorCode = "RATE_LIMIT_EXCEEDED"
)

// String implementa fmt.Stringer
func (e ErrorCode) String() string {
	return string(e)
}

// HTTPStatus mapea códigos de error a códigos HTTP apropiados
// Proporciona un mapping por defecto, pero los handlers pueden sobrescribirlo
func (e ErrorCode) HTTPStatus() int {
	switch e {
	// 401 Unauthorized
	case CodeUnauthorized, CodeTokenExpired, CodeTokenInvalid, CodePolicyVersionMismatch:
		return 401

	// 403 Forbidden
	case CodePermissionDenied, CodeInsufficientPermissions, CodeOperationNotAllowed:
		return 403

	// 404 Not Found
	case CodeNotFound, CodeMembershipNotFound:
		return 404

	// 409 Conflict
	case CodeAlreadyExists, CodeConflict, CodeResourceLocked:
		return 409

	// 410 Gone
	case CodeGone:
		return 410

	// 400 Bad Request
	case CodeValidationError, CodeInvalidRequest, CodeBadRequest,
		CodeMissingRequiredField, CodeInvalidFormat, CodeOutOfRange:
		return 400

	// 429 Too Many Requests
	case CodeRateLimitExceeded:
		return 429

	// 503 Service Unavailable
	case CodeServiceUnavailable:
		return 503

	// 504 Gateway Timeout
	case CodeTimeout:
		return 504

	// 500 Internal Server Error (default)
	default:
		return 500
	}
}
