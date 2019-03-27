package bus

import (
	"github.com/rafaeljesus/nsq-event-bus"
	"github.com/spf13/viper"
	"log"
	"os"
	"github.com/farzadrastegar/simple-cab/gateway"
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
	//set parameters from yaml
	nsqdLookupAddress := viper.GetString("urls.driverLocations.nsq.nsqLookupdAddress")
	nsqdAddress := viper.GetString("urls.driverLocations.nsq.nsqdAddress")
	topic := viper.GetString("urls.driverLocations.nsq.topic")
	channel := viper.GetString("urls.driverLocations.nsq.channel")
	maxInFlight := viper.GetInt("urls.driverLocations.nsq.maxInFlight")
	handlerConcurrency := viper.GetInt("urls.driverLocations.nsq.handlerConcurrency")

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
