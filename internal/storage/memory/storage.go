package memorystorage

import (
	"sync"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
)

type Storage struct {
	IPSubnets map[string]common.IPSubnet
	mu        sync.RWMutex
}

func New() common.StorageDriverInterface {
	return &Storage{
		IPSubnets: make(map[string]common.IPSubnet, 0),
	}
}

func (s *Storage) Add(ipSubnet common.IPSubnet) (common.IPSubnet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.IPSubnets[ipSubnet.ID.(string)]; exists {
		return ipSubnet, common.ErrIPSubnetAlreadyExists
	}

	s.IPSubnets[ipSubnet.ID.(string)] = ipSubnet
	return ipSubnet, nil
}

func (s *Storage) Update(ipSubnet common.IPSubnet) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, exist := s.IPSubnets[ipSubnet.ID.(string)]
	if !exist {
		return common.ErrIPSubnetNotFound
	}

	s.IPSubnets[ipSubnet.ID.(string)] = ipSubnet
	return nil
}

func (s *Storage) Delete(id interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.IPSubnets[id.(string)]; !exists {
		return common.ErrIPSubnetNotFound
	}

	delete(s.IPSubnets, id.(string))
	return nil
}

func (s *Storage) GetByID(id interface{}) (common.IPSubnet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	IPSubnet, ok := s.IPSubnets[id.(string)]
	if !ok {
		return common.IPSubnet{}, common.ErrIPSubnetNotFound
	}
	return IPSubnet, nil
}

func (s *Storage) List() ([]*common.IPSubnet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]*common.IPSubnet, 0, len(s.IPSubnets))
	for _, v := range s.IPSubnets {
		result = append(result, &v)
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
