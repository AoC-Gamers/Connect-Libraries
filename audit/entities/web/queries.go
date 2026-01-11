package web

// BuildSelectQuery construye la query base para seleccionar entradas de web audit
// Nota: scope_id es opcional para web audit
func BuildSelectQuery() string {
	return `SELECT id, scope_id, action, performed_by, payload, created_at 
	        FROM ` + TableName + ` 
	        WHERE 1=1`
}

// BuildCountQuery construye la query para contar entradas de web audit
func BuildCountQuery() string {
	return `SELECT COUNT(*) FROM ` + TableName + ` WHERE 1=1`
}

// BuildInsertQuery construye la query para insertar una entrada de web audit
func BuildInsertQuery() string {
	return `INSERT INTO ` + TableName + ` (scope_id, action, performed_by, payload, created_at) 
	        VALUES ($1, $2, $3, $4, $5)`
}

// BuildDeleteQuery construye la query para eliminar entradas antiguas (mantenimiento)
func BuildDeleteQuery() string {
	return `DELETE FROM ` + TableName + ` WHERE created_at < $1`
}
