//revive:disable
package common

import (
	"context"
	"time"
)

type IPSubnet struct {
	ID         interface{} `db:"id"`
	Subnet     string      `db:"sub_net"`
	DateCreate time.Time   `db:"date_create"`
}

func NewIPSubnet(subnet string, dateCreate time.Time) *IPSubnet {
	return &IPSubnet{
		Subnet:     subnet,
		DateCreate: dateCreate,
	}
}

type StorageDriverInterface interface {
	Add(IPSubnet IPSubnet) (IPSubnet, error)
	Update(IPSubnet IPSubnet) error
	Delete(id interface{}) error
	GetByID(id interface{}) (IPSubnet, error)
	List() ([]*IPSubnet, error)
	PrepareStorage(log LoggerInterface) error
	IsOverlapping(IPSubnet *IPSubnet) (bool, error)
}

type Storage struct {
	s   StorageDriverInterface
	ctx *context.Context
}

func New(ctx *context.Context, s StorageDriverInterface) (*Storage, error) {
	return &Storage{
		s:   s,
		ctx: ctx,
	}, nil
}

func (s *Storage) Add(ipSubnet IPSubnet) (IPSubnet, error) {
	return s.s.Add(ipSubnet)
}

func (s *Storage) Update(ipSubnet IPSubnet) error {
	return s.s.Update(ipSubnet)
}

func (s *Storage) Delete(id interface{}) error {
	return s.s.Delete(id)
}

func (s *Storage) GetByID(id interface{}) (IPSubnet, error) {
	return s.s.GetByID(id)
}

func (s *Storage) List() ([]*IPSubnet, error) {
	return s.s.List()
}

func (s *Storage) IsOverlapping(ipSubnet *IPSubnet) (bool, error) {
	return s.s.IsOverlapping(ipSubnet)
}
