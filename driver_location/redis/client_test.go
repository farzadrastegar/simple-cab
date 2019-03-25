package redis_test

import (
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/farzadrastegar/simple-cab/driver_location/redis"
	go_redis "github.com/go-redis/redis"
)

// Client represents a test wrapper to redis.Client.
type Client struct {
	*redis.Client

	db *go_redis.Client

	//// Services.
	//cabService driver_location.CabService
}

// NewClient returns a new instance of Client.
func NewClient() *Client {
	//driver_location.SetConfigFilename("../cmd/config.yaml")

	c := &Client{
		Client: redis.NewClient(),
	}
	c.db = c.GetDB()
	//c.cabService = c.Connect()

	return c
}

// MustOpenClient returns an new, open instance of Client.
func MustOpenClient() *Client {
	c := NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	return c
}

// Close closes the client and removes the underlying database.
func (c *Client) Close() error {
	if c != nil && c.db != nil {
		return c.db.Close()
	}
	return nil
}
