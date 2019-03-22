package mock

import (
	"github.com/farzadrastegar/simple-cab/zombie_driver"
)

type CabService struct {
	CheckZombieStatusFn      func(id string) (*zombie_driver.Status, error)
	CheckZombieStatusInvoked bool
}

func (s *CabService) CheckZombieStatus(id string) (*zombie_driver.Status, error) {
	s.CheckZombieStatusInvoked = true
	return s.CheckZombieStatusFn(id)
}
