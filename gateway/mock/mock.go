package mock

import (
	"github.com/farzadrastegar/simple-cab/gateway"
)

type CabService struct {
	StoreLocationFn      func(id string, data *gateway.Data) error
	StoreLocationInvoked bool

	CheckZombieStatusFn      func(id string) (*gateway.Status, error)
	CheckZombieStatusInvoked bool
}

func (s *CabService) StoreLocation(id string, data *gateway.Data) error {
	s.StoreLocationInvoked = true
	return s.StoreLocationFn(id, data)
}

func (s *CabService) CheckZombieStatus(id string) (*gateway.Status, error) {
	s.CheckZombieStatusInvoked = true
	return s.CheckZombieStatusFn(id)
}

type Bus struct {
	InitializeFn      func() (gateway.BusService, error)
	InitializeInvoked bool
}

func (b *Bus) Initialize() (gateway.BusService, error) {
	b.InitializeInvoked = true
	return b.InitializeFn()
}

type BusService struct {
	ProduceFn      func(id gateway.DriverID, message *gateway.Data) error
	ProduceInvoked bool
}

func (s *BusService) Produce(id gateway.DriverID, message *gateway.Data) error {
	s.ProduceInvoked = true
	return s.ProduceFn(id, message)
}

//// DefaultBusService returns input parameters in string format as its error
//func DefaultBusService() BusService {
//	return BusService{
//		ProduceFn: func(id gateway.DriverID, message *gateway.Data) error {
//			str := fmt.Sprintf("id=%d,{latitude:%f,longitude:%f}", id, message.Latitude, message.Longitude)
//
//			return gateway.Error(str)
//		},
//	}
//}
