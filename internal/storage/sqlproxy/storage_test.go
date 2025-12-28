//go:build integrations

package sqlproxy

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
)

func TestSqlProxyStorage(t *testing.T) {
	ipSubnet1 := common.NewIPSubnet("127.0.0.0/24", time.Now())
	ipSubnet2 := common.NewIPSubnet("10.0.0.0/24", ipSubnet1.DateCreate.Add(time.Second*-1))

	c, err := config.New("./test/config.yaml")

	require.NoError(t, err)
	ctx := context.Background()
	storageBlackList := New(&ctx, c.Storage)
	storageWhiteList := New(&ctx, c.Storage)

	t.Run("Add Subnet BlackList", func(t *testing.T) {
		newIPSubnet, err := storageBlackList.Add("black_list", *ipSubnet1)
		require.NoError(t, err)
		require.Equal(t, ipSubnet1, newIPSubnet)

		newIPSubnet2, err := storageBlackList.Add("black_list", *ipSubnet2)
		require.NoError(t, err)
		require.Equal(t, ipSubnet2, newIPSubnet2)
	})

	t.Run("Delete Subnet BlackList", func(t *testing.T) {
		err = storageBlackList.Delete("black_list", ipSubnet1.Subnet)
		require.NoError(t, err)

		_, err = storageBlackList.Get("black_list", ipSubnet1.Subnet)
		require.ErrorIs(t, common.ErrIPSubnetNotFound, err)
	})

	t.Run("Add Subnet WhiteList", func(t *testing.T) {
		newIPSubnet, err := storageWhiteList.Add("white_list", *ipSubnet1)
		require.NoError(t, err)
		require.Equal(t, ipSubnet1, newIPSubnet)

		newIPSubnet2, err := storageWhiteList.Add("white_list", *ipSubnet2)
		require.NoError(t, err)
		require.Equal(t, ipSubnet2, newIPSubnet2)
	})

	t.Run("Delete Subnet WhiteList", func(t *testing.T) {
		err = storageWhiteList.Delete("white_list", ipSubnet1.Subnet)
		require.NoError(t, err)

		_, err = storageWhiteList.Get("white_list", ipSubnet1.Subnet)
		require.ErrorIs(t, common.ErrIPSubnetNotFound, err)
	})

	t.Cleanup(func() {
		storageBlackList.Clear("white_list")
		storageWhiteList.Clear("black_list")
	})
}
