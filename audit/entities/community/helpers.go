package community

import (
	"encoding/json"
	"fmt"
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

// FormatCreatedPayload formatea el payload para acción CREATED
func FormatCreatedPayload(name, status string, hasOwner bool) string {
	return fmt.Sprintf(`{"name":"%s","status":"%s","hasOwner":%t}`, name, status, hasOwner)
}

// FormatStatusChangePayload formatea el payload para cambios de estado
func FormatStatusChangePayload(oldStatus, newStatus string) string {
	return fmt.Sprintf(`{"oldStatus":"%s","newStatus":"%s"}`, oldStatus, newStatus)
}

// FormatServerPayload formatea el payload para operaciones de servidor
func FormatServerPayload(serverID int64, serverName string) string {
	return fmt.Sprintf(`{"serverId":%d,"name":"%s"}`, serverID, serverName)
}

// FormatConfigPayload formatea el payload para cambios de configuración
func FormatConfigPayload(previousMode, newMode string, listSize int) string {
	return fmt.Sprintf(`{"previousMode":"%s","newMode":"%s","listSize":%d}`, previousMode, newMode, listSize)
}
