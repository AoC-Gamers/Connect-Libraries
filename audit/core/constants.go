package core

import "errors"

// Constantes SQL para construcción de queries
const (
	// Fragmentos de filtros
	SQLFilterAction       = " AND action = $%d"
	SQLFilterPerformedBy  = " AND performed_by = $%d"
	SQLFilterCreatedAtGTE = " AND created_at >= $%d"
	SQLFilterCreatedAtLTE = " AND created_at <= $%d"
	SQLFilterScopeID      = " AND scope_id = $%d"

	// Ordenamiento y paginación
	SQLOrderByCreatedDesc = " ORDER BY created_at DESC"
	SQLOrderByCreatedAsc  = " ORDER BY created_at ASC"
	SQLLimit              = " LIMIT $%d"
	SQLOffset             = " OFFSET $%d"
)

// Errores comunes
var (
	ErrEntryNil         = errors.New("audit entry is nil")
	ErrInvalidLimit     = errors.New("invalid limit: must be >= 0")
	ErrInvalidOffset    = errors.New("invalid offset: must be >= 0")
	ErrInvalidDateRange = errors.New("invalid date range: start date must be before end date")
	ErrEmptyAction      = errors.New("action cannot be empty")
	ErrEmptyPerformedBy = errors.New("performed_by cannot be empty")
)
