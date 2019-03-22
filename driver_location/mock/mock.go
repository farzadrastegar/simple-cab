package mock

import (
	"github.com/farzadrastegar/simple-cab/driver_location"
)

type CabService struct {
	StoreLocationFn      func(id string, loc *driver_location.Location) error
	StoreLocationInvoked bool

	GetDriverLocationsFn      func(id string, minutes float64) (*driver_location.Locations, error)
	GetDriverLocationsInvoked bool
}

func (s *CabService) StoreLocation(id string, loc *driver_location.Location) error {
	s.StoreLocationInvoked = true
	return s.StoreLocationFn(id, loc)
}

func (s *CabService) GetDriverLocations(id string, minutes float64) (*driver_location.Locations, error) {
	s.GetDriverLocationsInvoked = true
	return s.GetDriverLocationsFn(id, minutes)
}

type Bus struct {
	InitializeFn      func(cs CabService) (driver_location.BusService, error)
	InitializeInvoked bool
}

func (b *Bus) Initialize(cs CabService) (driver_location.BusService, error) {
	b.InitializeInvoked = true
	return b.InitializeFn(cs)
}

type BusService struct {
	ConsumeFn      func() error
	ConsumeInvoked bool
}

func (s *BusService) Consume() error {
	s.ConsumeInvoked = true
	return s.ConsumeFn()
}

//// Create a default id and location.
//const defaultMsgId = 123
//
//var defaultLocation = driver_location.Location{Latitude: 1.234567, Longitude: 1.234567}
//var DefaultBusServiceOut string
//
//// DefaultBusService returns input parameters in string format as its error.
//func DefaultBusService() BusService {
//	return BusService{
//		ConsumeFn: func() error {
//			DefaultBusServiceOut = fmt.Sprintf("id=%d,{latitude:%f,longitude:%f}", defaultMsgId, defaultLocation.Latitude, defaultLocation.Longitude)
//
//			return nil
//		},
//	}
//}
