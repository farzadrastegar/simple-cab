package zombie_driver

// DriverID represents a driver identifier.
type DriverID int

// Status represents driver status based on records in database.
type Status struct {
	ID     DriverID `json:"id"`
	Zombie bool     `json:"zombie"`
}

// Client creates a connection to the services.
type Client interface {
	Connect() CabService
}

// CabService represents a service for managing requests.
type CabService interface {
	CheckZombieStatus(id string) (*Status, error)
}

