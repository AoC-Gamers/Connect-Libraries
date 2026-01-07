package endpoints

// Connect-Core internal API endpoints
// Base URL: /core/internal

const (
	PathCoreMissionByID = "/core/internal/missions/{id}"
)

var (
	// User Management
	CoreSyncUser = Endpoint{
		Path:        "/core/internal/users/{steamid}",
		Method:      "POST",
		Description: "Sync user from Connect-Auth after login",
		RequiresKey: true,
		UsedBy:      []string{"auth"},
	}

	// Mission Management
	CoreGetMission = Endpoint{
		Path:        PathCoreMissionByID,
		Method:      "GET",
		Description: "Get mission by ID",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	CoreCreateMission = Endpoint{
		Path:        "/core/internal/missions",
		Method:      "POST",
		Description: "Create a new mission",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	CoreUpdateMission = Endpoint{
		Path:        PathCoreMissionByID,
		Method:      "PUT",
		Description: "Update a mission",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	CoreDeleteMission = Endpoint{
		Path:        PathCoreMissionByID,
		Method:      "DELETE",
		Description: "Delete a mission",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	// Gamemode Management
	CoreListGamemodes = Endpoint{
		Path:        "/core/internal/gamemodes",
		Method:      "GET",
		Description: "List all gamemodes",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	CoreCreateGamemode = Endpoint{
		Path:        "/core/internal/gamemodes",
		Method:      "POST",
		Description: "Create a new gamemode",
		RequiresKey: true,
		UsedBy:      []string{"auth"},
	}

	// Team Management
	CoreGetTeamMembers = Endpoint{
		Path:        "/core/internal/teams/{id}/members",
		Method:      "GET",
		Description: "Get team members",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	// Server Management
	CoreListServers = Endpoint{
		Path:        "/core/internal/servers",
		Method:      "GET",
		Description: "List all servers",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	CoreGetServer = Endpoint{
		Path:        "/core/internal/servers/{id}",
		Method:      "GET",
		Description: "Get server by ID",
		RequiresKey: true,
		UsedBy:      []string{"auth", "lobby"},
	}

	// Settings Management
	CoreGetSettings = Endpoint{
		Path:        "/core/internal/settings",
		Method:      "GET",
		Description: "Get global settings",
		RequiresKey: true,
		UsedBy:      []string{"auth", "rt"},
	}

	CoreUpdateSettings = Endpoint{
		Path:        "/core/internal/settings",
		Method:      "PUT",
		Description: "Update global settings",
		RequiresKey: true,
		UsedBy:      []string{"auth"},
	}
)

// CoreEndpoints returns all Connect-Core internal endpoints
func CoreEndpoints() []Endpoint {
	return []Endpoint{
		CoreSyncUser,
		CoreGetMission,
		CoreCreateMission,
		CoreUpdateMission,
		CoreDeleteMission,
		CoreListGamemodes,
		CoreCreateGamemode,
		CoreGetTeamMembers,
		CoreListServers,
		CoreGetServer,
		CoreGetSettings,
		CoreUpdateSettings,
	}
}

// CoreService represents the Connect-Core service
var CoreService = Service{
	Name:      "Connect-Core",
	BaseURL:   "/core/internal",
	Endpoints: CoreEndpoints(),
}
