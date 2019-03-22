package bus

import (
	"github.com/farzadrastegar/simple-cab/gateway"
	"github.com/farzadrastegar/simple-cab/gateway/config"
	"github.com/rafaeljesus/nsq-event-bus"
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
var _ gateway.Bus = &Bus{}

// Bus represents an implementation of gateway.Bus.
type Bus struct {
	Params    *Params
	ParamsSet bool
	Logger    *log.Logger
}

func (b *Bus) readParams() error {
	//create a config handler
	configHandler := config.NewConfig(b.Logger)

	//read yaml config
	configHandler.ReadYaml(gateway.GetConfigFilename())
	yaml := configHandler.GetYamlValue("urls", "driverLocations", "nsq")

	//set parameters from yaml
	msg := "ERROR: reading NSQ parameters failed"
	nsqdLookupAddress, err := yaml.Get("nsqLookupdAddress").String()
	if err != nil {
		return gateway.Error(msg)
	}
	nsqdAddress, err := yaml.Get("nsqdAddress").String()
	if err != nil {
		return gateway.Error(msg)
	}
	topic, err := yaml.Get("topic").String()
	if err != nil {
		return gateway.Error(msg)
	}
	channel, err := yaml.Get("channel").String()
	if err != nil {
		return gateway.Error(msg)
	}
	maxInFlight, err := yaml.Get("maxInFlight").Int()
	if err != nil {
		return gateway.Error(msg)
	}
	handlerConcurrency, err := yaml.Get("handlerConcurrency").Int()
	if err != nil {
		return gateway.Error(msg)
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

// Initialize initializes an NSQ bus service.
func (b *Bus) Initialize() (gateway.BusService, error) {
	if !b.ParamsSet {
		err := b.readParams()
		if err != nil {
			return nil, err
		}
	}

	return b.newBusService()
}

// newBusService creates a new bus service which implements gateway.BusService.
func (b *Bus) newBusService() (gateway.BusService, error) {
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
		emitter: emitter,
		logger:  b.Logger,
		topic:   b.Params.Topic,
	}, nil
}

// NewBus creates a new bus.
func NewBus() *Bus {
	return &Bus{
		ParamsSet: false,
		Logger:    log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile),
	}
}
