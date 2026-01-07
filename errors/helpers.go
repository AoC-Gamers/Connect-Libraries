package errors

import (
	"fmt"
	"net/http"
)

// ============================================
// Authentication & Authorization Helpers
// ============================================

// RespondUnauthorized responde con error 401 genérico
func RespondUnauthorized(w http.ResponseWriter, detail string) {
	RespondWithDetail(w, http.StatusUnauthorized, CodeUnauthorized,
		"unauthorized", detail)
}

// RespondUnauthorizedSimple responde con 401 sin detalle adicional
func RespondUnauthorizedSimple(w http.ResponseWriter) {
	RespondSimpleError(w, http.StatusUnauthorized, CodeUnauthorized,
		"unauthorized")
}

// RespondTokenExpired responde cuando el JWT ha expirado
func RespondTokenExpired(w http.ResponseWriter) {
	RespondError(w, http.StatusUnauthorized, CodeTokenExpired,
		"unauthorized",
		"JWT token has expired",
		map[string]interface{}{
			"should_refresh": true,
		})
}

// RespondTokenInvalid responde cuando el JWT es inválido
func RespondTokenInvalid(w http.ResponseWriter, reason string) {
	RespondError(w, http.StatusUnauthorized, CodeTokenInvalid,
		"unauthorized",
		fmt.Sprintf("JWT token is invalid: %s", reason),
		map[string]interface{}{
			"should_reauth": true,
		})
}

// RespondPolicyVersionMismatch responde cuando hay mismatch de policy version
func RespondPolicyVersionMismatch(w http.ResponseWriter, tokenVersion, currentVersion int) {
	RespondError(w, http.StatusUnauthorized, CodePolicyVersionMismatch,
		"policy version mismatch",
		"Token policy version does not match current version",
		map[string]interface{}{
			"token_version":   tokenVersion,
			"current_version": currentVersion,
			"should_reauth":   true,
		})
}

// RespondPermissionDenied responde cuando falta un permiso específico
func RespondPermissionDenied(w http.ResponseWriter, scopeID int64, scopeType, permission string, hasMembership bool) {
	detail := fmt.Sprintf("User lacks %s permission in %s scope %d", permission, scopeType, scopeID)
	if !hasMembership {
		detail = fmt.Sprintf("User does not have membership in %s scope %d", scopeType, scopeID)
	}

	RespondError(w, http.StatusForbidden, CodePermissionDenied,
		"insufficient permissions",
		detail,
		map[string]interface{}{
			"scope_id":            scopeID,
			"scope_type":          scopeType,
			"required_permission": permission,
			"has_membership":      hasMembership,
		})
}

// RespondInsufficientPermissions responde con permisos insuficientes genérico
func RespondInsufficientPermissions(w http.ResponseWriter, action string) {
	RespondWithDetail(w, http.StatusForbidden, CodeInsufficientPermissions,
		"insufficient permissions",
		fmt.Sprintf("You don't have permission to %s", action))
}

// ============================================
// Resource Error Helpers
// ============================================

// RespondNotFound responde cuando un recurso no existe
func RespondNotFound(w http.ResponseWriter, resourceType, identifier string) {
	RespondError(w, http.StatusNotFound, CodeNotFound,
		"not found",
		fmt.Sprintf("%s '%s' not found", resourceType, identifier),
		map[string]interface{}{
			"resource_type": resourceType,
			"identifier":    identifier,
		})
}

// RespondMembershipNotFound responde cuando no existe membership
func RespondMembershipNotFound(w http.ResponseWriter, scopeID int64, scopeType string, userID int64) {
	RespondError(w, http.StatusNotFound, CodeMembershipNotFound,
		"membership not found",
		fmt.Sprintf("User %d does not have membership in %s scope %d", userID, scopeType, scopeID),
		map[string]interface{}{
			"scope_id":   scopeID,
			"scope_type": scopeType,
			"user_id":    userID,
		})
}

// RespondNotFoundWithDetail responde con not found con solo un mensaje de detalle
func RespondNotFoundWithDetail(w http.ResponseWriter, detail string) {
	RespondError(w, http.StatusNotFound, CodeNotFound,
		"not found", detail, nil)
}

// RespondAlreadyExists responde cuando un recurso ya existe
func RespondAlreadyExists(w http.ResponseWriter, resourceType, identifier string) {
	RespondError(w, http.StatusConflict, CodeAlreadyExists,
		"resource already exists",
		fmt.Sprintf("%s '%s' already exists", resourceType, identifier),
		map[string]interface{}{
			"resource_type": resourceType,
			"identifier":    identifier,
		})
}

// RespondConflict responde con conflicto de estado
func RespondConflict(w http.ResponseWriter, reason string) {
	RespondWithDetail(w, http.StatusConflict, CodeConflict,
		"conflict", reason)
}

// ============================================
// Validation Error Helpers
// ============================================

// RespondValidationError responde con error de validación
func RespondValidationError(w http.ResponseWriter, field, reason string) {
	RespondError(w, http.StatusBadRequest, CodeValidationError,
		"validation failed",
		fmt.Sprintf("Field '%s': %s", field, reason),
		map[string]interface{}{
			"field":  field,
			"reason": reason,
		})
}

// RespondValidationErrors responde con múltiples errores de validación
func RespondValidationErrors(w http.ResponseWriter, errors map[string]string) {
	RespondError(w, http.StatusBadRequest, CodeValidationError,
		"validation failed",
		"Multiple validation errors occurred",
		map[string]interface{}{
			"errors": errors,
		})
}

// RespondValidationErrorWithDetail responde con error de validación con solo un mensaje de detalle
func RespondValidationErrorWithDetail(w http.ResponseWriter, detail string) {
	RespondError(w, http.StatusBadRequest, CodeValidationError,
		"validation failed", detail, nil)
}

// RespondBadRequest responde con bad request genérico
func RespondBadRequest(w http.ResponseWriter, detail string) {
	RespondWithDetail(w, http.StatusBadRequest, CodeBadRequest,
		"bad request", detail)
}

// RespondMissingField responde cuando falta un campo requerido
func RespondMissingField(w http.ResponseWriter, fieldName string) {
	RespondError(w, http.StatusBadRequest, CodeMissingRequiredField,
		"missing required field",
		fmt.Sprintf("Field '%s' is required", fieldName),
		map[string]interface{}{
			"field": fieldName,
		})
}

// RespondInvalidFormat responde cuando el formato es inválido
func RespondInvalidFormat(w http.ResponseWriter, field, expectedFormat string) {
	RespondError(w, http.StatusBadRequest, CodeInvalidFormat,
		"invalid format",
		fmt.Sprintf("Field '%s' has invalid format. Expected: %s", field, expectedFormat),
		map[string]interface{}{
			"field":           field,
			"expected_format": expectedFormat,
		})
}

// RespondOutOfRange responde cuando un valor está fuera de rango
func RespondOutOfRange(w http.ResponseWriter, field string, min, max, provided interface{}) {
	RespondError(w, http.StatusBadRequest, CodeOutOfRange,
		"value out of range",
		fmt.Sprintf("Field '%s' must be between %v and %v, got %v", field, min, max, provided),
		map[string]interface{}{
			"field":    field,
			"min":      min,
			"max":      max,
			"provided": provided,
		})
}

// ============================================
// Server Error Helpers
// ============================================

// RespondInternalError responde con error interno del servidor
func RespondInternalError(w http.ResponseWriter, detail string) {
	RespondWithDetail(w, http.StatusInternalServerError, CodeInternalError,
		"internal server error", detail)
}

// RespondInternalErrorSimple responde con error interno sin detalle adicional
func RespondInternalErrorSimple(w http.ResponseWriter) {
	RespondSimpleError(w, http.StatusInternalServerError, CodeInternalError,
		"internal server error")
}

// RespondDatabaseError responde con error de base de datos desde un error
func RespondDatabaseError(w http.ResponseWriter, err error) {
	detail := "Database operation failed"
	if err != nil {
		detail = fmt.Sprintf("Database operation failed: %s", err.Error())
	}
	RespondError(w, http.StatusInternalServerError, CodeDatabaseError,
		"database error", detail, nil)
}

// RespondDatabaseErrorWithOperation responde con error de base de datos con operación específica
func RespondDatabaseErrorWithOperation(w http.ResponseWriter, operation string) {
	RespondError(w, http.StatusInternalServerError, CodeDatabaseError,
		"database error",
		fmt.Sprintf("Database operation failed: %s", operation),
		map[string]interface{}{
			"operation": operation,
		})
}

// RespondServiceUnavailable responde cuando el servicio no está disponible
func RespondServiceUnavailable(w http.ResponseWriter, service string) {
	RespondError(w, http.StatusServiceUnavailable, CodeServiceUnavailable,
		"service unavailable",
		fmt.Sprintf("Service '%s' is temporarily unavailable", service),
		map[string]interface{}{
			"service":     service,
			"retry_after": 60, // segundos
		})
}

// RespondTimeout responde cuando una operación excede el timeout
func RespondTimeout(w http.ResponseWriter, operation string, timeoutSeconds int) {
	RespondError(w, http.StatusGatewayTimeout, CodeTimeout,
		"operation timeout",
		fmt.Sprintf("Operation '%s' exceeded timeout of %d seconds", operation, timeoutSeconds),
		map[string]interface{}{
			"operation":       operation,
			"timeout_seconds": timeoutSeconds,
		})
}

// ============================================
// Business Logic Error Helpers
// ============================================

// RespondOperationNotAllowed responde cuando una operación no está permitida
func RespondOperationNotAllowed(w http.ResponseWriter, operation, reason string) {
	RespondError(w, http.StatusForbidden, CodeOperationNotAllowed,
		"operation not allowed",
		fmt.Sprintf("Operation '%s' is not allowed: %s", operation, reason),
		map[string]interface{}{
			"operation": operation,
			"reason":    reason,
		})
}

// RespondQuotaExceeded responde cuando se excede una cuota
func RespondQuotaExceeded(w http.ResponseWriter, quotaType string, limit, current int) {
	RespondError(w, http.StatusTooManyRequests, CodeQuotaExceeded,
		"quota exceeded",
		fmt.Sprintf("Quota for '%s' exceeded. Limit: %d, Current: %d", quotaType, limit, current),
		map[string]interface{}{
			"quota_type": quotaType,
			"limit":      limit,
			"current":    current,
		})
}

// RespondRateLimitExceeded responde cuando se excede el rate limit
func RespondRateLimitExceeded(w http.ResponseWriter, limit int, window string, retryAfter int) {
	RespondError(w, http.StatusTooManyRequests, CodeRateLimitExceeded,
		"rate limit exceeded",
		fmt.Sprintf("Rate limit of %d requests per %s exceeded", limit, window),
		map[string]interface{}{
			"limit":       limit,
			"window":      window,
			"retry_after": retryAfter,
		})
}
