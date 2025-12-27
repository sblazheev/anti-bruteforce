//go:build integrations

package sqlstorage

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common"
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
)

func TestSqlStorage(t *testing.T) {
	ipSubnet1 := common.NewIPSubnet("127.0.0.0/24", time.Now())
	ipSubnet2 := common.NewIPSubnet("127.0.0.0/16", ipSubnet1.DateCreate.Add(time.Second*-1))
	ipSubnet3 := common.NewIPSubnet("127.0.0.0/8", ipSubnet1.DateCreate.Add(time.Second*10))
	ipSubnet4 := common.NewIPSubnet("127.0.0.128/25", ipSubnet1.DateCreate.Add(time.Second*5))

	c, err := config.New("./test/config.yaml")
	require.NoError(t, err)
	ctx := context.Background()
	s := New(&ctx, c.Storage)

	t.Run("Add subnet", func(t *testing.T) {
		new, err := s.Add("white_list", *ipSubnet1)
		require.NoError(t, err)
		require.Equal(t, ipSubnet1, new)

		new, err = s.Add("white_list", *ipSubnet2)
		require.NoError(t, err)
		require.Equal(t, ipSubnet2, new)

		new, err = s.Add("white_list", *ipSubnet3)
		require.NoError(t, err)
		require.Equal(t, ipSubnet3, new)
	})

	t.Run("Add IsOverlapping", func(t *testing.T) {
		IsOverlapping, err := s.IsOverlapping("white_list", ipSubnet4)
		require.NoError(t, err)
		require.Equal(t, true, IsOverlapping)
	})

	t.Run("Get list subnet", func(t *testing.T) {
		list, err := s.List("white_list")
		require.NoError(t, err)
		require.NotNil(t, list)
	})
	t.Run("Get subnet", func(t *testing.T) {
		new, err := s.Get("white_list", ipSubnet1.Subnet)
		require.NoError(t, err)
		require.Equal(t, ipSubnet1.Subnet, new.Subnet)
	})
	t.Run("In subnet", func(t *testing.T) {
		inSubNet, err := s.InSubNet("white_list", "127.0.0.1")
		require.NoError(t, err)
		require.Equal(t, true, inSubNet)
	})
	t.Run("Delete subnet", func(t *testing.T) {
		err := s.Delete("white_list", ipSubnet1.Subnet)
		require.NoError(t, err)
		_, err = s.Get("white_list", ipSubnet1.Subnet)
		require.Equal(t, err, common.ErrIPSubnetNotFound)
	})

	t.Cleanup(func() {
		//s.Clear("white_list")
	})
}
