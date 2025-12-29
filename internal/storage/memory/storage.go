package memorystorage

import (
	"net/netip"
	"sync"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
	"go4.org/netipx"                                      //nolint:depguard
)

type Storage struct {
	IPSubnets     map[string]*common.IPSubnet
	mu            sync.RWMutex
	builderNetipx netipx.IPSetBuilder
	ipSet         *netipx.IPSet
}

func New() common.StorageDriverInterface {
	return &Storage{
		IPSubnets: make(map[string]*common.IPSubnet, 0),
	}
}

func (s *Storage) Add(_ string, ipSubnet common.IPSubnet) (*common.IPSubnet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.IPSubnets[ipSubnet.Subnet]; exists {
		return &ipSubnet, nil
	}

	s.IPSubnets[ipSubnet.Subnet] = &ipSubnet

	net, err := netip.ParsePrefix(ipSubnet.Subnet)
	if err != nil {
		return nil, err
	}
	s.builderNetipx.AddPrefix(net)

	s.ipSet, err = s.builderNetipx.IPSet()
	if err != nil {
		return nil, err
	}

	return &ipSubnet, nil
}

func (s *Storage) Update(_ string, ipSubnet common.IPSubnet) error {
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

	delete(s.IPSubnets, subnet)
	return nil
}

func (s *Storage) Clear(_ string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	clear(s.IPSubnets)
	return nil
}

func (s *Storage) Get(_ string, subnet string) (*common.IPSubnet, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	IPSubnet, ok := s.IPSubnets[subnet]
	if !ok {
		return &common.IPSubnet{}, common.ErrIPSubnetNotFound
	}
	return IPSubnet, nil
}

func (s *Storage) List(_ string) ([]*common.IPSubnet, error) {
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

func (s *Storage) IsOverlapping(_ string, _ *common.IPSubnet) (bool, error) {
	return false, nil
}

func (s *Storage) InSubNet(_ string, ip string) (bool, error) {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return false, err
	}
	if s.ipSet == nil {
		return false, nil
	}
	return s.ipSet.Contains(addr), nil
}

func (s *Storage) Load(_ string) (bool, error) {
	return true, nil
}
