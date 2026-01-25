package endpoints

// Connect-Auth internal API endpoints
// Base URL: /auth/internal

var (
	// Permission Management
	AuthCheckPermission = Endpoint{
		Path:        "/auth/internal/permissions/check",
		Method:      "POST",
		Description: "Check if user has a specific permission",
		RequiresKey: true,
		UsedBy:      []string{"core", "rt", "lobby"},
	}

	AuthGetCatalog = Endpoint{
		Path:        "/auth/internal/catalog",
		Method:      "GET",
		Description: "Get complete permissions catalog",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	// Web Owner
	AuthGetWebOwner = Endpoint{
		Path:        "/auth/internal/owner",
		Method:      "GET",
		Description: "Get web_owner (global owner) steamID",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	// Scope Management
	AuthCreateScope = Endpoint{
		Path:        "/auth/internal/scopes",
		Method:      "POST",
		Description: "Create a new scope (community/team/lobby)",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthDeleteScope = Endpoint{
		Path:        "/auth/internal/scopes/{scopeId}",
		Method:      "DELETE",
		Description: "Delete a scope and all related data",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthGetScopeStaffCount = Endpoint{
		Path:        "/auth/internal/scopes/{scopeType}/{entityId}/staff-count",
		Method:      "GET",
		Description: "Get staff count for a scope",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	// Role Management
	AuthAssignRole = Endpoint{
		Path:        "/auth/internal/roles/assign",
		Method:      "POST",
		Description: "Assign a role to a user in a scope",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthRemoveRole = Endpoint{
		Path:        "/auth/internal/roles/{userId}/{scopeId}",
		Method:      "DELETE",
		Description: "Remove a role from a user in a scope",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	// Membership Management
	AuthCreateMembership = Endpoint{
		Path:        "/auth/internal/memberships",
		Method:      "POST",
		Description: "Create a membership for a user",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthCheckMembership = Endpoint{
		Path:        "/auth/internal/memberships/check",
		Method:      "POST",
		Description: "Check if user is member of a scope",
		RequiresKey: true,
		UsedBy:      []string{"core", "rt"},
	}

	AuthUpdateMembership = Endpoint{
		Path:        "/auth/internal/memberships/{userId}/{scopeId}",
		Method:      "PATCH",
		Description: "Update membership role and/or permissions",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	// Ownership Transfer
	AuthInitiateTransfer = Endpoint{
		Path:        "/auth/internal/scopes/{scopeType}/{entityId}/transfers/initiate",
		Method:      "POST",
		Description: "Initiate ownership transfer",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthGetPendingTransfer = Endpoint{
		Path:        "/auth/internal/scopes/{scopeType}/{entityId}/transfers/pending",
		Method:      "GET",
		Description: "Get pending ownership transfer",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthRespondToTransfer = Endpoint{
		Path:        "/auth/internal/transfers/{transferId}/respond",
		Method:      "POST",
		Description: "Respond to ownership transfer (accept/reject)",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthCancelTransfer = Endpoint{
		Path:        "/auth/internal/transfers/{transferId}/cancel",
		Method:      "POST",
		Description: "Cancel ownership transfer",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthGetTransfer = Endpoint{
		Path:        "/auth/internal/transfers/{transferId}",
		Method:      "GET",
		Description: "Get ownership transfer by ID",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	// Cache Management
	AuthCacheRefresh = Endpoint{
		Path:        "/auth/internal/cache/refresh",
		Method:      "POST",
		Description: "Publish cache invalidation event",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthCacheReload = Endpoint{
		Path:        "/auth/internal/cache/reload",
		Method:      "POST",
		Description: "Reload cache completely",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	AuthCacheStatus = Endpoint{
		Path:        "/auth/internal/cache/status",
		Method:      "GET",
		Description: "Get cache status",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	// Notifications
	AuthCreateNotification = Endpoint{
		Path:        "/auth/internal/notifications",
		Method:      "POST",
		Description: "Create a notification for a user",
		RequiresKey: true,
		UsedBy:      []string{"rt", "lobby"},
	}

	AuthTeamMatchLobbyNotification = Endpoint{
		Path:        "/auth/internal/notifications/team-match-lobby",
		Method:      "POST",
		Description: "Notify team members about lobby match",
		RequiresKey: true,
		UsedBy:      []string{"rt", "lobby"},
	}
)

// AuthEndpoints returns all Connect-Auth internal endpoints
func AuthEndpoints() []Endpoint {
	return []Endpoint{
		AuthCheckPermission,
		AuthGetCatalog,
		AuthGetWebOwner,
		AuthCreateScope,
		AuthDeleteScope,
		AuthGetScopeStaffCount,
		AuthAssignRole,
		AuthRemoveRole,
		AuthCreateMembership,
		AuthCheckMembership,
		AuthUpdateMembership,
		AuthInitiateTransfer,
		AuthGetPendingTransfer,
		AuthRespondToTransfer,
		AuthCancelTransfer,
		AuthGetTransfer,
		AuthCacheRefresh,
		AuthCacheReload,
		AuthCacheStatus,
		AuthCreateNotification,
		AuthTeamMatchLobbyNotification,
	}
}

// AuthService represents the Connect-Auth service
var AuthService = Service{
	Name:      "Connect-Auth",
	BaseURL:   "/auth/internal",
	Endpoints: AuthEndpoints(),
}
