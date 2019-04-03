package bus

import (
	"fmt"

	"github.com/farzadrastegar/simple-cab/gateway"

	bus "github.com/rafaeljesus/nsq-event-bus"
	logger "github.com/sirupsen/logrus"
)

type Event struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Id        string  `json:"id"`
}

//type Event struct {
//	body string
//	id   gateway.DriverID
//}

// Ensure service implements interface.
var _ gateway.BusService = &BusService{}

// BusService represents an implementation of gateway.BusService.
type BusService struct {
	emitter *bus.Emitter
	topic   string
}

// Produce writes a message into bus.
func (s *BusService) Produce(id gateway.DriverID, message *gateway.Data) error {
	e := Event{Latitude: message.Latitude, Longitude: message.Longitude, Id: fmt.Sprintf("%d", int(id))}

	//send message via NSQ.
	if err := s.emitter.EmitAsync(s.topic, &e); err != nil {
		logger.Printf("ERROR: publishing NSQ message failed %#v", err)
		return err
	} else {
		logger.Printf("[Message published] %#v", e)
	}
	return nil
}
