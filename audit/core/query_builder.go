package core

import "fmt"

// ApplyFilters aplica los filtros comunes a una query SQL
// Retorna la query modificada y los argumentos actualizados
func (f *Filters) ApplyFilters(baseQuery string, args []interface{}) (string, []interface{}) {
	query := baseQuery

	if f.Action != "" {
		query += fmt.Sprintf(SQLFilterAction, len(args)+1)
		args = append(args, f.Action)
	}

	if f.PerformedBy != "" {
		query += fmt.Sprintf(SQLFilterPerformedBy, len(args)+1)
		args = append(args, f.PerformedBy)
	}

	if f.StartDate != nil {
		query += fmt.Sprintf(SQLFilterCreatedAtGTE, len(args)+1)
		args = append(args, f.StartDate)
	}

	if f.EndDate != nil {
		query += fmt.Sprintf(SQLFilterCreatedAtLTE, len(args)+1)
		args = append(args, f.EndDate)
	}

	return query, args
}

// ApplyPagination aplica ordenamiento y paginación a una query SQL
// Retorna la query modificada y los argumentos actualizados
func (f *Filters) ApplyPagination(query string, args []interface{}) (string, []interface{}) {
	return f.ApplyPaginationWithOrder(query, args, SQLOrderByCreatedDesc)
}

// ApplyPaginationWithOrder aplica ordenamiento personalizado y paginación
func (f *Filters) ApplyPaginationWithOrder(query string, args []interface{}, orderBy string) (string, []interface{}) {
	query += orderBy

	if f.Limit > 0 {
		query += fmt.Sprintf(SQLLimit, len(args)+1)
		args = append(args, f.Limit)
	}

	if f.Offset > 0 {
		query += fmt.Sprintf(SQLOffset, len(args)+1)
		args = append(args, f.Offset)
	}

	return query, args
}

// ApplyScopeIDFilter aplica filtro de scope_id opcional (para web audit)
func (f *Filters) ApplyScopeIDFilter(query string, args []interface{}) (string, []interface{}) {
	if f.ScopeID > 0 {
		query += fmt.Sprintf(SQLFilterScopeID, len(args)+1)
		args = append(args, f.ScopeID)
	}
	return query, args
}
