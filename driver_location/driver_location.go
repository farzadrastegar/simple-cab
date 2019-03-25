package driver_location

// DriverID represents a driver identifier.
type DriverID int

// Location represents geo information from cars.
type Location struct {
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	Updated_at string  `json:"updated_at,omitempty"`
}

// Locations represents records in database related to a car.
type Locations struct {
	Locations []Location `json:"locations,omitempty"`
	Err       string     `json:"err,omitempty"`
	ServerIP  string     `json:"serverIP,omitempty"`
}

// Bus represents a connection to a message bus.
type Bus interface {
	Initialize(cs CabService) (BusService, error)
}

// BusService represents a service for handling a bus.
type BusService interface {
	Consume() error
}

// Client creates a connection to the services.
type Client interface {
	Connect() CabService
}

// CabService represents a service for managing requests.
type CabService interface {
	StoreLocation(id string, loc *Location) error
	GetDriverLocations(id string, minutes float64) (*Locations, error)
}

