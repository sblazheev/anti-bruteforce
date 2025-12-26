package memorystorage

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
)

func TestMemoryStorage(t *testing.T) {
	ipSubnet1 := common.NewIPSubnet("127.0.0.0/24", time.Now())

	ipSubnet2 := common.NewIPSubnet("127.0.0.0/16", ipSubnet1.DateCreate.Add(time.Second*-1))

	t.Run("Storage create", func(t *testing.T) {
		s := New()
		require.Equal(t, &Storage{IPSubnets: make(map[string]common.IPSubnet, 0)}, s)
	})
	t.Run("Add IP Subnet", func(t *testing.T) {
		storage := New()

		newIPSubnet, err := storage.Add(*ipSubnet1)
		require.NoError(t, err)
		ipSubnet1.ID = newIPSubnet.ID
		require.Equal(t, ipSubnet1, &newIPSubnet)

		newIPSubnet2, err := storage.Add(*ipSubnet2)
		require.NoError(t, err)
		ipSubnet2.ID = newIPSubnet2.ID
		require.Equal(t, ipSubnet2, &newIPSubnet2)
	})

	t.Run("Delete event", func(t *testing.T) {
		storage := New()

		newIPSubNet, err := storage.Add(*ipSubnet1)
		require.NoError(t, err)
		require.Equal(t, ipSubnet1, &newIPSubNet)

		err = storage.Delete(newIPSubNet.ID)
		require.NoError(t, err)

		_, err = storage.GetByID(newIPSubNet.ID)
		require.ErrorIs(t, common.ErrEventNotFound, err)
	})
}
