package sqlproxy

import (
	"context"

	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common"                       //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config"                       //nolint:depguard
	memorystorage "gitlab.wsrubi.ru/go/anti-bruteforce/internal/storage/memory" //nolint:depguard
	sqlstorage "gitlab.wsrubi.ru/go/anti-bruteforce/internal/storage/sql"       //nolint:depguard
)

type Storage struct {
	memory     common.StorageDriverInterface
	sqlstorage common.StorageDriverInterface
	ctx        *context.Context
}

func New(ctx *context.Context, c config.StorageConfig) common.StorageDriverInterface {
	storage := &Storage{
		ctx:        ctx,
		memory:     memorystorage.New(),
		sqlstorage: sqlstorage.New(ctx, c),
	}
	return storage
}

func (s *Storage) Add(jar string, ipSubnet common.IPSubnet) (*common.IPSubnet, error) {
	return s.memory.Add(jar, ipSubnet)
}

func (s *Storage) Update(jar string, ipSubnet common.IPSubnet) error {
	return s.memory.Update(jar, ipSubnet)
}

func (s *Storage) Delete(jar string, subnet string) error {
	return s.memory.Delete(jar, subnet)
}

func (s *Storage) Clear(jar string) error {
	s.memory.Clear(jar)
	return nil
}

func (s *Storage) Get(jar string, subnet string) (*common.IPSubnet, error) {
	return s.memory.Get(jar, subnet)
}

func (s *Storage) List(jar string) ([]*common.IPSubnet, error) {
	return s.sqlstorage.List(jar)
}

func (s *Storage) PrepareStorage(_ common.LoggerInterface) error {
	return nil
}

func (s *Storage) IsOverlapping(_ string, _ *common.IPSubnet) (bool, error) {
	return false, nil
}

func (s *Storage) InSubNet(jar string, ip string) (bool, error) {
	return s.memory.InSubNet(jar, ip)
}

func (s *Storage) Load(jar string) (bool, error) {
	ipSubNet, err := s.List(jar)
	if err != nil {
		return false, err
	}

	for _, itemSubNet := range ipSubNet {
		s.memory.Add(jar, *itemSubNet)
	}

	return true, nil
}
