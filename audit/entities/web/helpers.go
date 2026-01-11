package web

import (
	"encoding/json"
	"fmt"

	"github.com/AoC-Gamers/Connect-Libraries/audit/core"
)

// FormatPayloadJSON convierte un map a JSON string para el payload
func FormatPayloadJSON(data map[string]interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}
	return string(bytes), nil
}

// FormatSimplePayload crea un payload simple con un campo
func FormatSimplePayload(key, value string) string {
	return fmt.Sprintf(`{"%s":"%s"}`, key, value)
}

// ApplyWebFilters aplica filtros específicos para web audit
// Web audit permite scope_id opcional (para eventos globales del sistema)
func ApplyWebFilters(filters *core.Filters, baseQuery string, args []interface{}) (string, []interface{}) {
	query := baseQuery

	// Aplicar scope_id solo si está presente
	query, args = filters.ApplyScopeIDFilter(query, args)

	// Aplicar filtros comunes
	query, args = filters.ApplyFilters(query, args)

	return query, args
}

// FormatLoginPayload formatea el payload para login/logout
func FormatLoginPayload(ipAddress, userAgent string) string {
	return fmt.Sprintf(`{"ip":"%s","userAgent":"%s"}`, ipAddress, userAgent)
}

// FormatPermissionPayload formatea el payload para cambios de permisos
func FormatPermissionPayload(permission, scope string, scopeID int64) string {
	return fmt.Sprintf(`{"permission":"%s","scope":"%s","scopeId":%d}`, permission, scope, scopeID)
}

// FormatRolePayload formatea el payload para cambios de rol
func FormatRolePayload(role, targetUser string) string {
	return fmt.Sprintf(`{"role":"%s","targetUser":"%s"}`, role, targetUser)
}

// EmptyPayload retorna un payload vacío
func EmptyPayload() string {
	return "{}"
}
