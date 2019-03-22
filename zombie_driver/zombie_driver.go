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

// CreateCabService creates a service through input client.
func CreateCabService(c Client) CabService {
	return c.Connect()
}

// zombieDriverConfigFilename contains configuration parameters in yaml format.
var zombieDriverConfigFilename string

// SetYamlFilename sets the configurations file name.
func SetConfigFilename(yamlFilename string) {
	zombieDriverConfigFilename = yamlFilename
}

// SetYamlFilename sets the configurations file name.
func GetConfigFilename() string {
	return zombieDriverConfigFilename
}
