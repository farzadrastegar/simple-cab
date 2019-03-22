package bus_test

import (
	"github.com/farzadrastegar/simple-cab/gateway/bus"
)

var testParams = &bus.Params{
	NSQLookupdAddress:  ":4161",
	NSQDAddress:        ":4150",
	Topic:              "test_locations",
	MaxInFlight:        25,
	Channel:            "test_store_location",
	HandlerConcurrency: 4,
}

// Bus is a test wrapper for bus.Bus.
type Bus struct {
	*bus.Bus
}

// NewBus creates a new bus wrapper instance.
func NewBus() *Bus {

	// Create Bus wrapper.
	b := &Bus{
		Bus: bus.NewBus(),
	}

	// Set parameters.
	b.Params = testParams
	b.ParamsSet = true

	return b
}
