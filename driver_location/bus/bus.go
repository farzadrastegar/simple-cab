package bus

import (
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/farzadrastegar/simple-cab/driver_location/config"
	"github.com/rafaeljesus/nsq-event-bus"
	"github.com/spf13/viper"
	"log"
	"os"
)

//const YamlFilename = "../config.yaml"

type Params struct {
	NSQLookupdAddress  string
	NSQDAddress        string
	Topic              string
	MaxInFlight        int
	Channel            string
	HandlerConcurrency int
}

// Ensure service implements interface.
var _ driver_location.Bus = &Bus{}

// Bus represents an implementation of driver_location.Bus.
type Bus struct {
	Params    *Params
	ParamsSet bool
	Logger    *log.Logger
	//BusService *BusService
}

func (b *Bus) readParams0() error {
	//create a config handler
	configHandler := config.NewConfig(b.Logger)

	//read yaml config
	configHandler.ReadYaml(driver_location.GetConfigFilename())
	yaml := configHandler.GetYamlValue("urls", "driverLocations", "nsq")

	//set parameters from yaml
	msg := "ERROR: reading NSQ parameters failed"
	nsqdLookupAddress, err := yaml.Get("nsqLookupdAddress").String()
	if err != nil {
		return driver_location.Error(msg)
	}
	nsqdAddress, err := yaml.Get("nsqdAddress").String()
	if err != nil {
		return driver_location.Error(msg)
	}
	topic, err := yaml.Get("topic").String()
	if err != nil {
		return driver_location.Error(msg)
	}
	channel, err := yaml.Get("channel").String()
	if err != nil {
		return driver_location.Error(msg)
	}
	maxInFlight, err := yaml.Get("maxInFlight").Int()
	if err != nil {
		return driver_location.Error(msg)
	}
	handlerConcurrency, err := yaml.Get("handlerConcurrency").Int()
	if err != nil {
		return driver_location.Error(msg)
	}

	b.ParamsSet = true

	b.Params = &Params{
		NSQLookupdAddress:  nsqdLookupAddress,
		NSQDAddress:        nsqdAddress,
		Topic:              topic,
		MaxInFlight:        maxInFlight,
		Channel:            channel,
		HandlerConcurrency: handlerConcurrency,
	}
	return nil
}

func (b *Bus) readParams() error {
	//set parameters from yaml
	nsqdLookupAddress = viper.GetString("urls.driverLocations.nsq.nsqLookupdAddress")
	nsqdAddress = viper.GetString("urls.driverLocations.nsq.nsqdAddress")
	topic = viper.GetString("urls.driverLocations.nsq.topic")
	channel = viper.GetString("urls.driverLocations.nsq.channel")
	maxInFlight = viper.GetInt("urls.driverLocations.nsq.maxInFlight")
	handlerConcurrency = viper.GetInt("urls.driverLocations.nsq.handlerConcurrency")

	b.ParamsSet = true

	b.Params = &Params{
		NSQLookupdAddress:  nsqdLookupAddress,
		NSQDAddress:        nsqdAddress,
		Topic:              topic,
		MaxInFlight:        maxInFlight,
		Channel:            channel,
		HandlerConcurrency: handlerConcurrency,
	}
	return nil
}

// Initialize initializes an NSQ bus service.
func (b *Bus) Initialize(cs driver_location.CabService) (driver_location.BusService, error) {
	// Read parameters.
	if !b.ParamsSet {
		err := b.readParams()
		if err != nil {
			return nil, err
		}
	}

	// Initialize an emitter for the bus.
	emitter, err := bus.NewEmitter(bus.EmitterConfig{
		Address:     b.Params.NSQDAddress,
		MaxInFlight: b.Params.MaxInFlight,
	})
	if err != nil {
		return nil, err
	}

	return &BusService{
		emitter:    emitter,
		logger:     b.Logger,
		nsqData:    b.Params,
		cabService: cs,
	}, nil
}

// NewBus creates a new bus.
func NewBus() *Bus {
	// Create a bus instance.
	b := &Bus{
		ParamsSet: false,
		Logger:    log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
	}

	return b
}
