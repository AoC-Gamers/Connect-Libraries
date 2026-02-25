package authjwt

// GetSteamID retorna el SteamID del usuario desde claims
func (c *Claims) GetSteamID() string {
	return c.SteamID
}

// GetRole retorna el rol del usuario
func (c *Claims) GetRole() string {
	return c.Role
}

// HasPermission verifica si el usuario tiene un permiso específico (bitmask)
func (c *Claims) HasPermission(permission uint64) bool {
	// Primero verificar si está explícitamente denegado
	if (c.DenyPermissions & permission) != 0 {
		return false
	}
	// Luego verificar si está permitido
	return (c.AllowPermissions & permission) != 0
}

// IsAdmin verifica si el usuario es administrador (web_admin o web_owner)
func (c *Claims) IsAdmin() bool {
	return c.Role == "web_admin" || c.Role == "web_owner"
}

// IsOwner verifica si el usuario es owner
func (c *Claims) IsOwner() bool {
	return c.Role == "web_owner"
}

// IsModerator verifica si el usuario es moderador o superior
func (c *Claims) IsModerator() bool {
	return c.Role == "web_moderator" || c.IsAdmin()
}

// IsStaff verifica si el usuario es staff o superior (incluye moderator, admin, owner)
func (c *Claims) IsStaff() bool {
	return c.Role == "web_staff" || c.IsModerator()
}
