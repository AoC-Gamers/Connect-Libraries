package team

import "fmt"

const (
	// TableName es el nombre de la tabla de auditoría de equipos
	TableName = "audit.team_audit"

	// Actions - Acciones de auditoría para equipos
	ActionCreated          = "TEAM_CREATED"
	ActionUpdated          = "TEAM_UPDATED"
	ActionDeleted          = "TEAM_DELETED"
	ActionOwnerTransferred = "TEAM_OWNER_TRANSFERRED"
	ActionSettingsUpdated  = "TEAM_SETTINGS_UPDATED"
	ActionSuspended        = "TEAM_SUSPENDED"
	ActionActivated        = "TEAM_ACTIVATED"
)

// validActions contiene todas las acciones válidas para team audit
var validActions = map[string]bool{
	ActionCreated:          true,
	ActionUpdated:          true,
	ActionDeleted:          true,
	ActionOwnerTransferred: true,
	ActionSettingsUpdated:  true,
	ActionSuspended:        true,
	ActionActivated:        true,
}

// ValidateAction valida que una acción sea válida para team audit
func ValidateAction(action string) error {
	if !validActions[action] {
		return fmt.Errorf("invalid team action: %s", action)
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
