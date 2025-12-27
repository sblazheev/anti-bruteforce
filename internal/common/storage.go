//revive:disable
package common

import (
	"context"
	"time"
)

type IPSubnet struct {
	ID         int64     `db:"id"`
	Subnet     string    `db:"sub_net"`
	DateCreate time.Time `db:"date_create"`
}

func NewIPSubnet(subnet string, dateCreate time.Time) *IPSubnet {
	return &IPSubnet{
		Subnet:     subnet,
		DateCreate: dateCreate,
	}
}

type StorageDriverInterface interface {
	Add(jar string, ipSubnet IPSubnet) (*IPSubnet, error)
	Update(jar string, ipSubnet IPSubnet) error
	Delete(jar string, subnet string) error
	Get(jar string, subnet string) (*IPSubnet, error)
	List(jar string) ([]*IPSubnet, error)
	Clear(jar string) error
	PrepareStorage(log LoggerInterface) error
	IsOverlapping(ipSubnet *IPSubnet) (bool, error)
}

type Storage struct {
	jar string
	s   StorageDriverInterface
	ctx *context.Context
}

func NewStorage(jar string, ctx *context.Context, s StorageDriverInterface) (*Storage, error) {
	return &Storage{
		s:   s,
		ctx: ctx,
	}, nil
}

func (s *Storage) Add(ipSubnet IPSubnet) (*IPSubnet, error) {
	return s.s.Add(s.jar, ipSubnet)
}

func (s *Storage) Update(ipSubnet IPSubnet) error {
	return s.s.Update(s.jar, ipSubnet)
}

func (s *Storage) Delete(subnet string) error {
	return s.s.Delete(s.jar, subnet)
}

func (s *Storage) Get(subnet string) (*IPSubnet, error) {
	return s.s.Get(s.jar, subnet)
}

func (s *Storage) List() ([]*IPSubnet, error) {
	return s.s.List(s.jar)
}

func (s *Storage) Clear() error {
	return s.s.Clear(s.jar)
}

func (s *Storage) IsOverlapping(ipSubnet *IPSubnet) (bool, error) {
	return s.s.IsOverlapping(ipSubnet)
}
