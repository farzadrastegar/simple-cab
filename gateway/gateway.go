package gateway

// DriverID represents a driver identifier.
type DriverID int

// Data represents geo information from cars.
type Data struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Status represents driver status based on records in database.
type Status struct {
	ID     DriverID `json:"id"`
	Zombie bool     `json:"zombie"`
}

// Bus represents a connection to a message bus.
type Bus interface {
	Initialize() (BusService, error)
}

// BusService represents a service for handling a bus.
type BusService interface {
	Produce(id DriverID, message *Data) error
}

// Client creates a connection to the services.
type Client interface {
	Connect() CabService
}

// CabService represents a service for managing data.
type CabService interface {
	StoreLocation(id string, data *Data) error
	CheckZombieStatus(id string) (*Status, error)
}
