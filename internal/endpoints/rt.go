package endpoints

// Connect-RT internal API endpoints
// Base URL: /rt/internal

var (
	// Presence Management
	RTGetUserPresence = Endpoint{
		Path:        "/rt/internal/presence/{steamid}",
		Method:      "GET",
		Description: "Get user presence by steamID",
		RequiresKey: true,
		UsedBy:      []string{"core", "auth"},
	}

	RTInitializePresence = Endpoint{
		Path:        "/rt/internal/presence/{steamid}",
		Method:      "POST",
		Description: "Initialize user presence after login",
		RequiresKey: true,
		UsedBy:      []string{"auth"},
	}

	RTBatchGetPresence = Endpoint{
		Path:        "/rt/internal/presence/batch",
		Method:      "POST",
		Description: "Get presence for multiple users (max 100)",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	RTGetOnlineUsers = Endpoint{
		Path:        "/rt/internal/presence/online",
		Method:      "GET",
		Description: "Get all online users",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}

	RTGetUsersByStatus = Endpoint{
		Path:        "/rt/internal/presence/by-status/{status}",
		Method:      "GET",
		Description: "Get users by presence status",
		RequiresKey: true,
		UsedBy:      []string{"core"},
	}
)

// RTEndpoints returns all Connect-RT internal endpoints
func RTEndpoints() []Endpoint {
	return []Endpoint{
		RTGetUserPresence,
		RTInitializePresence,
		RTBatchGetPresence,
		RTGetOnlineUsers,
		RTGetUsersByStatus,
	}
}

// RTService represents the Connect-RT service
var RTService = Service{
	Name:      "Connect-RT",
	BaseURL:   "/rt/internal",
	Endpoints: RTEndpoints(),
}
