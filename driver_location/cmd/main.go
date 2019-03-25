package main

import (
	"flag"
	"fmt"
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/farzadrastegar/simple-cab/driver_location/bus"
	"github.com/farzadrastegar/simple-cab/driver_location/http"
	"github.com/farzadrastegar/simple-cab/driver_location/redis"
	"github.com/spf13/viper"
	"os"
	"os/signal"
	"simple-cab/driver_location/config"
)

// Initialize configurations.
func init() {
	profile := flag.String("profile", "test", "Environment profile, something similar to spring profiles")
	configServerUrl := flag.String("configServerUrl", "http://configserver:8888", "Address to config server")
	configBranch := flag.String("configBranch", "master", "git branch to fetch configuration from")

	flag.Parse()

	fmt.Println("Specified configBranch is " + *configBranch)

	viper.Set("profile", *profile)
	viper.Set("configServerUrl", *configServerUrl)
	viper.Set("configBranch", *configBranch)
}

func main() {
	//// Set yaml filename.
	//driver_location.SetConfigFilename("../config.yaml")

	// Load configurations.
	config.LoadConfigurationFromBranch(
		viper.GetString("configServerUrl"),
		driver_location.AppName,
		viper.GetString("profile"),
		viper.GetString("configBranch"))

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
	fmt.Println("Got signal:", sig)
}
