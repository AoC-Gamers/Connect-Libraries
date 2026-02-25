package web

import "fmt"

const (
	// TableName es el nombre de la tabla de auditoría web/sistema
	TableName = "audit.web_audit"

	// Actions - Acciones de auditoría para operaciones web/sistema
	ActionUserLogin           = "USER_LOGIN"
	ActionUserLogout          = "USER_LOGOUT"
	ActionUserRegistered      = "USER_REGISTERED"
	ActionPasswordChanged     = "PASSWORD_CHANGED"
	ActionEmailChanged        = "EMAIL_CHANGED"
	ActionSystemConfigUpdated = "SYSTEM_CONFIG_UPDATED"
	ActionPermissionGranted   = "PERMISSION_GRANTED"
	ActionPermissionRevoked   = "PERMISSION_REVOKED"
	ActionRoleAssigned        = "ROLE_ASSIGNED"
	ActionRoleRemoved         = "ROLE_REMOVED"
	ActionAPIKeyCreated       = "API_KEY_CREATED"
	ActionAPIKeyRevoked       = "API_KEY_REVOKED"
	ActionSecurityAlert       = "SECURITY_ALERT"
)

// validActions contiene todas las acciones válidas para web audit
var validActions = map[string]bool{
	ActionUserLogin:           true,
	ActionUserLogout:          true,
	ActionUserRegistered:      true,
	ActionPasswordChanged:     true,
	ActionEmailChanged:        true,
	ActionSystemConfigUpdated: true,
	ActionPermissionGranted:   true,
	ActionPermissionRevoked:   true,
	ActionRoleAssigned:        true,
	ActionRoleRemoved:         true,
	ActionAPIKeyCreated:       true,
	ActionAPIKeyRevoked:       true,
	ActionSecurityAlert:       true,
}

// ValidateAction valida que una acción sea válida para web audit
func ValidateAction(action string) error {
	if !validActions[action] {
		return fmt.Errorf("invalid web action: %s", action)
	}
	return nil
}

// IsValidAction verifica si una acción es válida sin retornar error
func IsValidAction(action string) bool {
	return validActions[action]
}

// GetAllActions retorna todas las acciones válidas
func GetAllActions() []string {
	actions := make([]string, 0, len(validActions))
	for action := range validActions {
		actions = append(actions, action)
	}
	return actions
}
