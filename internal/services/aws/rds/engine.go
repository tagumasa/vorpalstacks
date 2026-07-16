package rds

// Engine manages per-resource (cluster/instance) engine lifecycles.
// Implementations open isolated engines with dynamically allocated ports
// for each resource created via the RDS API.
type Engine interface {
	Open(region, resourceID string) (int, error)
	Close(resourceID string) error
	EngineType() string
}

// GetPorter returns the allocated port for a resource.
type GetPorter interface {
	GetPort(resourceID string) (int, error)
}
