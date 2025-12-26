//go:build integrations

package sqlstorage

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require" //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common"
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/config" //nolint:depguard
)

func TestSqlStorage(t *testing.T) {
	event := *common.NewEvent("", "Test", time.Now(), 60, "Test", 0, 0)
	event1 := *common.NewEvent("", "Test 1", event.DateTime.Add(time.Second*60), 60, "Test 1", 1, 0)
	event2 := *common.NewEvent("", "Test 2", event.DateTime.Add(time.Second*120), 60, "Test 2", 1, 0)
	event3 := *common.NewEvent("", "Test 3", event.DateTime.Add(time.Second*180), 60, "Test 3", 1, 0)
	event4 := *common.NewEvent("", "Test 4", event.DateTime.Add(time.Second*240), 60, "Test 4", 1, 0)
	event5 := *common.NewEvent("", "Test 5", event.DateTime.Add(time.Hour*24), 60, "Test 5", 1, 0)
	c, err := config.New("./test/config.yaml")
	require.NoError(t, err)
	ctx := context.Background()
	s := New(&ctx, c.Storage)

	tx, _ := s.(*Storage).db.BeginTx(ctx, nil)
	t.Run("Add event", func(t *testing.T) {
		newEvent, err := s.Add(event)
		require.NoError(t, err)
		event.ID = newEvent.ID
		require.Equal(t, event, newEvent)
	})

	t.Run("Get list", func(t *testing.T) {
		newEvents, err := s.List()
		require.NoError(t, err)
		require.NotNil(t, newEvents)
	})
	t.Run("Get event", func(t *testing.T) {
		newEvent, err := s.GetByID(event.ID)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s %d %d %d %s %s %d", event.ID, event.User, event.Duration,
			event.NotifyTime, event.Description, event.Title, event.DateTime.Unix()),
			fmt.Sprintf("%s %d %d %d %s %s %d", newEvent.ID, newEvent.User, newEvent.Duration,
				newEvent.NotifyTime, newEvent.Description, newEvent.Title, newEvent.DateTime.Unix()))
	})

	t.Run("Update event", func(t *testing.T) {
		event.User = 1
		err := s.Update(event)
		require.NoError(t, err)

		newEvent, err := s.GetByID(event.ID)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf("%s %d %d %d %s %s %d", event.ID, event.User, event.Duration,
			event.NotifyTime, event.Description, event.Title, event.DateTime.Unix()),
			fmt.Sprintf("%s %d %d %d %s %s %d", newEvent.ID, newEvent.User, newEvent.Duration,
				newEvent.NotifyTime, newEvent.Description, newEvent.Title, newEvent.DateTime.Unix()))
	})

	t.Run("Delete event", func(t *testing.T) {
		err := s.Delete(event.ID)
		require.NoError(t, err)
		_, err = s.GetByID(event.ID)
		require.Equal(t, err, common.ErrEventNotFound)
	})

	t.Cleanup(func() {
		tx.Rollback()
	})
}
