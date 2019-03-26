package bus

import (
	"github.com/farzadrastegar/simple-cab/driver_location"
	"github.com/rafaeljesus/nsq-event-bus"
	"log"
	"time"
)

const NsqMaxDeliveryAttempts = 5

type Event struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Id        string  `json:"id"`
}

// Ensure service implements interface.
var _ driver_location.BusService = &BusService{}

// BusService represents an implementation of driver_location.BusService.
type BusService struct {
	emitter    *bus.Emitter
	logger     *log.Logger
	nsqData    *Params
	cabService driver_location.CabService
}

// Consume consumes a message through a bus.
func (s *BusService) Consume() error {
	return bus.On(bus.ListenerConfig{
		Lookup:             []string{s.nsqData.NSQLookupdAddress},
		Topic:              s.nsqData.Topic,
		Channel:            s.nsqData.Channel,
		HandlerFunc:        s.nsqHandlerFunc,
		HandlerConcurrency: s.nsqData.HandlerConcurrency,
	})
}

func (s *BusService) nsqHandlerFunc(message *bus.Message) (reply interface{}, err error) {
	startTime := time.Now()
	defer s.logger.Printf("message consumption processed in %s\n", time.Now().Sub(startTime))

	e := Event{}
	if err = message.DecodePayload(&e); err != nil {
		message.Finish()
		return
	}

	s.logger.Println("[Message to be consumed]", e)

	if message.Attempts > NsqMaxDeliveryAttempts {
		message.Finish()
		return
	}

	err = s.cabService.StoreLocation(e.Id, &driver_location.Location{Latitude: e.Latitude, Longitude: e.Longitude})
	//err = s.redisHandler.StorePayload(&e)
	if err != nil {
		s.logger.Printf("requeuing message")
		message.Requeue(time.Millisecond)
		//message.Finish() //todo
		return
	}

	s.logger.Println("[Message consumed]", e)

	message.Finish()
	return
}
