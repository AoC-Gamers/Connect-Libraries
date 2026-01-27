package permissions

// HasPermission checks if a given bitmask has a specific permission
// This is the core authorization check used across all services
func HasPermission(mask, permission uint64) bool {
	return (mask & permission) == permission
}

// HasAnyPermission checks if a bitmask has at least one of the specified permissions
func HasAnyPermission(mask uint64, permissions ...uint64) bool {
	for _, perm := range permissions {
		if HasPermission(mask, perm) {
			return true
		}
	}
	return false
}

// HasAllPermissions checks if a bitmask has all of the specified permissions
func HasAllPermissions(mask uint64, permissions ...uint64) bool {
	for _, perm := range permissions {
		if !HasPermission(mask, perm) {
			return false
		}
	}
	return true
}

// ApplyDenyMask applies a deny mask to an allow mask
// Deny permissions take precedence over allow permissions
func ApplyDenyMask(allowMask, denyMask uint64) uint64 {
	return allowMask & ^denyMask
}

// GetEffectivePermissions calculates the effective permissions after applying deny mask
// This is the final permission check used in authorization
func GetEffectivePermissions(allowMask, denyMask uint64) uint64 {
	return ApplyDenyMask(allowMask, denyMask)
}

// CanPerformAction checks if a user can perform an action after considering deny mask
// This is the main function used by middleware
func CanPerformAction(allowMask, denyMask, requiredPermission uint64) bool {
	effectiveMask := GetEffectivePermissions(allowMask, denyMask)
	return HasPermission(effectiveMask, requiredPermission)
}

// CanPerformAnyAction checks if a user can perform at least one of the specified actions
func CanPerformAnyAction(allowMask, denyMask uint64, requiredPermissions ...uint64) bool {
	effectiveMask := GetEffectivePermissions(allowMask, denyMask)
	return HasAnyPermission(effectiveMask, requiredPermissions...)
}

// CanPerformAllActions checks if a user can perform all of the specified actions
func CanPerformAllActions(allowMask, denyMask uint64, requiredPermissions ...uint64) bool {
	effectiveMask := GetEffectivePermissions(allowMask, denyMask)
	return HasAllPermissions(effectiveMask, requiredPermissions...)
}
