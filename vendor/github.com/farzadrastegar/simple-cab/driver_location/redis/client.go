package redis

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/farzadrastegar/simple-cab/driver_location/config"
	"github.com/go-redis/redis"
	"log"
	"os"
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

	logger   *log.Logger
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
	var err error

	//read database config from config.yaml
	dbConf := config.NewConfig(c.logger)
	dbConf.ReadYaml(driver_location.GetConfigFilename())
	serverYaml := dbConf.GetYamlValue("database")
	//read address
	c.params.dbAddr, err = serverYaml.Get("address").String()
	if err != nil {
		c.logger.Fatalf("database configurations in yaml not readable => %#v", err)
	}
	//read port
	c.params.dbPort, err = serverYaml.Get("port").String()
	if err != nil {
		c.logger.Fatalf("database configurations in yaml not readable => %#v", err)
	}
	//read password
	c.params.dbPassword, err = serverYaml.Get("password").String()
	if err != nil {
		c.params.dbPassword = ""
		c.logger.Printf("using database default password")
	}
	//read db
	c.params.dbDB, err = serverYaml.Get("db").Int()
	if err != nil {
		c.params.dbDB = 0
		c.logger.Printf("using default redis database")
	}
}

func NewClient() *Client {
	c := &Client{
		logger: log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
		params: &Params{dbAddr: "localhost", dbPort: "6379", dbPassword: "", dbDB: 0},
	}
	c.cabService.db = &c.db
	c.cabService.logger = c.logger
	c.cabService.now = &c.Now
	c.cabService.interval = &c.Interval

	return c
}

// Open opens and initializes the Redis database.
func (c *Client) Open() error {
	// Read database parameters.
	c.readParams()

	dbAddrPort := fmt.Sprintf("%s:%s", c.params.dbAddr, c.params.dbPort)
	c.logger.Printf("database is connecting to %s", dbAddrPort)

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

	c.logger.Println("redis is ready")
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
