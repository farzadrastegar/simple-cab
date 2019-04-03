package redis

import (
	"fmt"

	"github.com/farzadrastegar/simple-cab/driver_location"

	"github.com/go-redis/redis"
	logger "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

//const YamlFilename = "../config.yaml"

type Params struct {
	dbAddr     string
	dbPort     string
	dbPassword string
	dbDB       int
}

// Ensure *Client implements driver_location.Client.
var _ driver_location.Client = &Client{}

// Client represents a client to the underlying Redis data store.
type Client struct {
	db *redis.Client

	params   *Params
	Now      int64
	Interval int64

	// Services.
	cabService CabService
}

func (c *Client) GetDB() *redis.Client {
	return c.db
}

func (c *Client) readParams() {
	//read address
	c.params.dbAddr = viper.GetString("database.address")

	//read port
	c.params.dbPort = viper.GetString("database.port")

	//read password
	c.params.dbPassword = viper.GetString("database.password")

	//read db
	c.params.dbDB = viper.GetInt("database.db")
}

func NewClient() *Client {
	c := &Client{
		params: &Params{dbAddr: "localhost", dbPort: "6379", dbPassword: "", dbDB: 0},
	}
	c.cabService.db = &c.db
	c.cabService.now = &c.Now
	c.cabService.interval = &c.Interval

	return c
}

// Open opens and initializes the Redis database.
func (c *Client) Open() error {
	// Read database parameters.
	c.readParams()

	dbAddrPort := fmt.Sprintf("%s:%s", c.params.dbAddr, c.params.dbPort)
	logger.Printf("database is connecting to %s", dbAddrPort)

	// Create connection.
	c.db = redis.NewClient(&redis.Options{
		Addr:     dbAddrPort,
		Password: c.params.dbPassword,
		DB:       c.params.dbDB,
	})

	// Test connection.
	_, err := c.db.Ping().Result()
	if err != nil {
		c.db = nil
		return err
	}

	logger.Println("redis is ready")
	return nil
}

// Close closes the underlying Redis database.
func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Connect returns a new session to the BoltDB database.
func (c *Client) Connect() driver_location.CabService {
	return &c.cabService
}
