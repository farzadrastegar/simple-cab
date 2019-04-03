package bus

import (
	"github.com/farzadrastegar/simple-cab/driver_location"

	bus "github.com/rafaeljesus/nsq-event-bus"
	"github.com/spf13/viper"
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
	//BusService *BusService
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
		nsqData:    b.Params,
		cabService: cs,
	}, nil
}

// NewBus creates a new bus.
func NewBus() *Bus {
	// Create a bus instance.
	b := &Bus{
		ParamsSet: false,
	}

	return b
}
