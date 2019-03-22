package main

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/driver_location"
	DL_bus "github.com/farzadrastegar/simple-cab/driver_location/bus"
	DL_http "github.com/farzadrastegar/simple-cab/driver_location/http"
	DL_redis "github.com/farzadrastegar/simple-cab/driver_location/redis"
	"github.com/farzadrastegar/simple-cab/gateway"
	GW_bus "github.com/farzadrastegar/simple-cab/gateway/bus"
	GW_client "github.com/farzadrastegar/simple-cab/gateway/client"
	GW_http "github.com/farzadrastegar/simple-cab/gateway/http"
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	ZD_client "github.com/farzadrastegar/simple-cab/zombie_driver/client"
	ZD_http "github.com/farzadrastegar/simple-cab/zombie_driver/http"
	"os"
	"os/signal"
)

func main() {
	// Get current directory.
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	// Open servers.
	gateway.SetConfigFilename(dir + "/cmd/config_gateway.yaml")
	_, gwSrv := OpenGatewayServer()
	defer gwSrv.Close()

	zombie_driver.SetConfigFilename(dir + "/cmd/config_zombie_driver.yaml")
	_, zdSrv := OpenZombieDriverServer()
	defer zdSrv.Close()

	driver_location.SetConfigFilename(dir + "/cmd/config_driver_location.yaml")
	_, dlSrv := OpenDriverLocationServer()
	defer dlSrv.Close()

	// Block until an OS interrupt is received.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	sig := <-ch
	fmt.Println("Got signal:", sig)
}

// OpenDriverLocationServer returns a working driver_location server and client.
func OpenDriverLocationServer() (*DL_redis.Client, *DL_http.Server) {
	// Create a cab service through redis.
	c := DL_redis.NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	cabService := driver_location.CreateCabService(c)

	// Create a bus consumer.
	b := DL_bus.NewBus()
	busService, err := b.Initialize(cabService)
	if err != nil {
		panic(err)
	}
	err = busService.Consume()
	if err != nil {
		panic(err)
	}

	// Attach cabService to HTTP handler.
	h := DL_http.NewDataHandler()
	h.CabService = cabService

	// Start an HTTP server.
	srv := DL_http.NewServer()
	srv.Handler = &DL_http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return c, srv
}

// OpenZombieDriverServer returns a working zombie_driver server.
func OpenZombieDriverServer() (*ZD_client.Client, *ZD_http.Server) {
	// Create a client for managing services.
	c := ZD_client.NewClient()

	// Create cab services.
	s := zombie_driver.CreateCabService(c)

	// Attach services to HTTP handler.
	h := ZD_http.NewDataHandler()
	h.CabService = s

	// Start an HTTP server.
	srv := ZD_http.NewServer()
	srv.Handler = &ZD_http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return c, srv
}

// OpenGatewayServer returns a working gateway server.
func OpenGatewayServer() (*GW_client.Client, *GW_http.Server) {
	// Create a bus service.
	b := GW_bus.NewBus()
	busService, err := b.Initialize()
	if err != nil {
		panic(err)
	}

	// Create a client and attach bus service to it.
	c := GW_client.NewClient()
	c.Handler = GW_client.NewHandler()
	c.Handler.BusService = busService

	// Create cab services.
	s := gateway.CreateCabService(c)

	// Attach services to HTTP handler.
	h := GW_http.NewDataHandler()
	h.CabService = s

	// Start an HTTP server.
	srv := GW_http.NewServer()
	srv.Handler = &GW_http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}

	return c, srv
}
