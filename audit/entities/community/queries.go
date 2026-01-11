package community

// BuildSelectQuery construye la query base para seleccionar entradas de community audit
func BuildSelectQuery() string {
	return `SELECT id, scope_id, action, performed_by, payload, created_at 
	        FROM ` + TableName + ` 
	        WHERE scope_id = $1`
}

// BuildCountQuery construye la query para contar entradas de community audit
func BuildCountQuery() string {
	return `SELECT COUNT(*) FROM ` + TableName + ` WHERE scope_id = $1`
}

// BuildInsertQuery construye la query para insertar una entrada de community audit
func BuildInsertQuery() string {
	return `INSERT INTO ` + TableName + ` (scope_id, action, performed_by, payload, created_at) 
	        VALUES ($1, $2, $3, $4, $5)`
}

// BuildDeleteQuery construye la query para eliminar entradas antiguas (mantenimiento)
func BuildDeleteQuery() string {
	return `DELETE FROM ` + TableName + ` WHERE created_at < $1`
}
