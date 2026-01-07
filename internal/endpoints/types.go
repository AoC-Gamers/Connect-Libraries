package endpoints

// Endpoint represents an internal API endpoint definition
type Endpoint struct {
	Path        string
	Method      string
	Description string
	RequiresKey bool
	UsedBy      []string
}

// Service represents a backend service with internal endpoints
type Service struct {
	Name      string
	BaseURL   string
	Endpoints []Endpoint
}
