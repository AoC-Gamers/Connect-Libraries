package core

import "time"

// Filters define los filtros comunes para consultar entradas de auditoría
type Filters struct {
	ScopeID     int64      // ID del scope (community, team, etc.)
	Action      string     // Filtro opcional por acción
	PerformedBy string     // Filtro opcional por usuario (SteamID)
	StartDate   *time.Time // Filtro opcional por fecha de inicio (inclusivo)
	EndDate     *time.Time // Filtro opcional por fecha de fin (inclusivo)
	Limit       int        // Límite de resultados para paginación
	Offset      int        // Offset para paginación
}

// Validate valida que los filtros sean correctos
func (f *Filters) Validate() error {
	if f.Limit < 0 {
		return ErrInvalidLimit
	}
	if f.Offset < 0 {
		return ErrInvalidOffset
	}
	if f.StartDate != nil && f.EndDate != nil && f.StartDate.After(*f.EndDate) {
		return ErrInvalidDateRange
	}
	return nil
}

// SetDefaults establece valores por defecto para la paginación
func (f *Filters) SetDefaults() {
	if f.Limit <= 0 || f.Limit > 100 {
		f.Limit = 50
	}
	if f.Offset < 0 {
		f.Offset = 0
	}
}
