package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/bus"
	"github.com/farzadrastegar/simple-cab/gateway/client"
	"github.com/farzadrastegar/simple-cab/gateway/http"

	logger "github.com/sirupsen/logrus"
)

func main() {
	// Load configurations from input flags (i.e. -configServerUrl, -profile, -configBranch).
	gateway.LoadConfigurationFromBranch()

	// Create a bus service.
	b := bus.NewBus()
	busService, err := b.Initialize()
	if err != nil {
		log.Fatal(err)
	}

	// Create a client and attach bus service to it.
	c := client.NewClient()
	c.Handler = client.NewHandler()
	c.Handler.BusService = busService

	// Create cab services.
	s := gateway.CreateCabService(c)

	// Attach services to HTTP handler.
	h := http.NewDataHandler()
	h.CabService = s

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
