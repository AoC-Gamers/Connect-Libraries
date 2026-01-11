package team

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
func FormatCreatedPayload(name, tag, status string) string {
	return fmt.Sprintf(`{"name":"%s","tag":"%s","status":"%s"}`, name, tag, status)
}

// FormatStatusChangePayload formatea el payload para cambios de estado
func FormatStatusChangePayload(oldStatus, newStatus string) string {
	return fmt.Sprintf(`{"oldStatus":"%s","newStatus":"%s"}`, oldStatus, newStatus)
}

// FormatOwnerTransferPayload formatea el payload para transferencia de propiedad
func FormatOwnerTransferPayload(newOwner string) string {
	return fmt.Sprintf(`{"newOwner":"%s"}`, newOwner)
}

// EmptyPayload retorna un payload vacío
func EmptyPayload() string {
	return "{}"
}
