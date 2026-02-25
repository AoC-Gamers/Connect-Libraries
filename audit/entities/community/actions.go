package community

import "fmt"

const (
	// TableName es el nombre de la tabla de auditoría de comunidades
	TableName = "audit.community_audit"

	// Actions - Acciones de auditoría para comunidades
	ActionCreated               = "COMMUNITY_CREATED"
	ActionUpdated               = "COMMUNITY_UPDATED"
	ActionDeleted               = "COMMUNITY_DELETED"
	ActionSuspended             = "SUSPEND"
	ActionActivated             = "ACTIVATE"
	ActionSettingsUpdated       = "COMMUNITY_SETTINGS_UPDATED"
	ActionOwnerTransferred      = "FORCE_TRANSFER"
	ActionServerAdded           = "SERVER_ADDED"
	ActionServerUpdated         = "SERVER_UPDATED"
	ActionServerRemoved         = "SERVER_REMOVED"
	ActionMissionConfigUpdated  = "MISSION_CONFIG_UPDATED"
	ActionGamemodeConfigUpdated = "GAMEMODE_CONFIG_UPDATED"
)

// validActions contiene todas las acciones válidas para community audit
var validActions = map[string]bool{
	ActionCreated:               true,
	ActionUpdated:               true,
	ActionDeleted:               true,
	ActionSuspended:             true,
	ActionActivated:             true,
	ActionSettingsUpdated:       true,
	ActionOwnerTransferred:      true,
	ActionServerAdded:           true,
	ActionServerUpdated:         true,
	ActionServerRemoved:         true,
	ActionMissionConfigUpdated:  true,
	ActionGamemodeConfigUpdated: true,
}

// ValidateAction valida que una acción sea válida para community audit
func ValidateAction(action string) error {
	if !validActions[action] {
		return fmt.Errorf("invalid community action: %s", action)
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
