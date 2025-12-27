package memorystorage

import (
	"sync"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
)

type Storage struct {
	IPSubnets map[string]*common.IPSubnet
	mu        sync.RWMutex
}

func New() common.StorageDriverInterface {
	return &Storage{
		IPSubnets: make(map[string]*common.IPSubnet, 0),
	}
}

func (s *Storage) Add(jar string, ipSubnet common.IPSubnet) (*common.IPSubnet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.IPSubnets[ipSubnet.Subnet]; exists {
		return &ipSubnet, common.ErrIPSubnetAlreadyExists
	}

	s.IPSubnets[ipSubnet.Subnet] = &ipSubnet
	return &ipSubnet, nil
}

func (s *Storage) Update(jar string, ipSubnet common.IPSubnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.IPSubnets[ipSubnet.Subnet]
	if !exist {
		return common.ErrIPSubnetNotFound
	}

	s.IPSubnets[ipSubnet.Subnet] = &ipSubnet
	return nil
}

func (s *Storage) Delete(_ string, subnet string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.IPSubnets[subnet]; !exists {
		return common.ErrIPSubnetNotFound
	}

	delete(s.IPSubnets, subnet)
	return nil
}

func (s *Storage) Clear(_ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.IPSubnets)
	return nil
}

func (s *Storage) Get(jar string, subnet string) (*common.IPSubnet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	IPSubnet, ok := s.IPSubnets[subnet]
	if !ok {
		return &common.IPSubnet{}, common.ErrIPSubnetNotFound
	}
	return IPSubnet, nil
}

func (s *Storage) List(jar string) ([]*common.IPSubnet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*common.IPSubnet, 0, len(s.IPSubnets))
	for _, v := range s.IPSubnets {
		result = append(result, v)
	}
	return result, nil
}

func (s *Storage) PrepareStorage(_ common.LoggerInterface) error {
	return nil
}

func (s *Storage) IsOverlapping(_ *common.IPSubnet) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return false, nil
}
