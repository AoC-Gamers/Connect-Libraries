package web

import "fmt"

const (
	// TableName es el nombre de la tabla de auditoría web/sistema
	TableName = "audit.web_audit"
	actionAPI = "API"
	actionKey = "KEY"

	// ActionUserLogin registra el inicio de sesión de un usuario.
	ActionUserLogin = "USER_LOGIN"
	// ActionUserLogout registra el cierre de sesión de un usuario.
	ActionUserLogout = "USER_LOGOUT"
	// ActionUserRegistered registra el alta de un usuario.
	ActionUserRegistered = "USER_REGISTERED"
	// ActionPasswordChanged registra el cambio de contraseña.
	ActionPasswordChanged = "PASSWORD_CHANGED"
	// ActionEmailChanged registra el cambio de correo electrónico.
	ActionEmailChanged = "EMAIL_CHANGED"
	// ActionSystemConfigUpdated registra una actualización de configuración del sistema.
	ActionSystemConfigUpdated = "SYSTEM_CONFIG_UPDATED"
	// ActionPermissionGranted registra la concesión de un permiso.
	ActionPermissionGranted = "PERMISSION_GRANTED"
	// ActionPermissionRevoked registra la revocación de un permiso.
	ActionPermissionRevoked = "PERMISSION_REVOKED"
	// ActionRoleAssigned registra la asignación de un rol.
	ActionRoleAssigned = "ROLE_ASSIGNED"
	// ActionRoleRemoved registra la eliminación de un rol.
	ActionRoleRemoved = "ROLE_REMOVED"
	// ActionAPIKeyCreated registra la creación de una API key.
	ActionAPIKeyCreated = actionAPI + "_" + actionKey + "_CREATED"
	// ActionAPIKeyRevoked registra la revocación de una API key.
	ActionAPIKeyRevoked = actionAPI + "_" + actionKey + "_REVOKED"
	// ActionSecurityAlert registra una alerta de seguridad.
	ActionSecurityAlert = "SECURITY_ALERT"
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
