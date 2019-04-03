package main

import (
	"os"
	"os/signal"

	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/farzadrastegar/simple-cab/driver_location/bus"
	"github.com/farzadrastegar/simple-cab/driver_location/http"
	"github.com/farzadrastegar/simple-cab/driver_location/redis"

	logger "github.com/sirupsen/logrus"
)

func main() {
	// Load configurations from input flags (i.e. -configServerUrl, -profile, -configBranch).
	driver_location.LoadConfigurationFromBranch()

	// Create a cab service through redis.
	c := redis.NewClient()
	if err := c.Open(); err != nil {
		panic(err)
	}
	cabService := driver_location.CreateCabService(c)

	// Create a bus consumer.
	b := bus.NewBus()
	busService, err := b.Initialize(cabService)
	if err != nil {
		panic(err)
	}
	err = busService.Consume()
	if err != nil {
		panic(err)
	}

	// Attach cabService to HTTP handler.
	h := http.NewDataHandler()
	h.CabService = cabService

	// Start an HTTP server.
	srv := http.NewServer()
	srv.Handler = &http.Handler{DataHandler: h}
	if err := srv.Open(); err != nil {
		panic(err)
	}
	defer srv.Close()

	// Block until an OS interrupt is received.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	sig := <-ch
	logger.Println("Got signal:", sig)
}
