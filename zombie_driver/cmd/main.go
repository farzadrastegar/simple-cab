package main

import (
	"fmt"
	"github.com/farzadrastegar/simple-cab/zombie_driver"
	"github.com/farzadrastegar/simple-cab/zombie_driver/client"
	"github.com/farzadrastegar/simple-cab/zombie_driver/http"
	"os"
	"os/signal"
)

func main() {
	// Set yaml filename.
	zombie_driver.SetConfigFilename("../config.yaml")

	// Create a client for managing services.
	c := client.NewClient()

	// Create cab services.
	s := zombie_driver.CreateCabService(c)

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
	fmt.Println("Got signal:", sig)
}
