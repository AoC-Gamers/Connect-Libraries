package core

import "time"

// ValidateEntryTime valida que el timestamp de una entrada sea válido
func ValidateEntryTime(createdAt time.Time) error {
	if createdAt.IsZero() {
		return nil // Se asignará automáticamente
	}
	// No permitir timestamps futuros
	if createdAt.After(time.Now().Add(time.Minute)) {
		return ErrInvalidDateRange
	}
	return nil
}

// ValidateAction valida que una acción no esté vacía
func ValidateAction(action string) error {
	if action == "" {
		return ErrEmptyAction
	}
	return nil
}

// ValidatePerformedBy valida que performed_by no esté vacío
func ValidatePerformedBy(performedBy string) error {
	if performedBy == "" {
		return ErrEmptyPerformedBy
	}
	return nil
}

// EnsureTimestamp asegura que el timestamp esté establecido
func EnsureTimestamp(createdAt *time.Time) {
	if createdAt.IsZero() {
		*createdAt = time.Now()
	}
}
