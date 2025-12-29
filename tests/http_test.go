//go:build integration

package tests

import (
	"bytes"
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"                  //nolint:depguard
	"github.com/stretchr/testify/require"                 //nolint:depguard
	"gitlab.wsrubi.ru/go/anti-bruteforce/internal/common" //nolint:depguard
)

const (
	httpURL = "http://anti_bruteforce:80"
)

func TestHTTPApi(t *testing.T) {
	ipSubnet1 := common.NewIPSubnet("127.0.0.0/24", time.Now())
	ipSubnet2 := common.NewIPSubnet("10.0.0.0/24", ipSubnet1.DateCreate.Add(time.Second*-1))

	t.Run("Create subnet black list", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/add/black-list",
			bytes.NewBuffer([]byte("{\"net\":\""+ipSubnet1.Subnet+"\"}")))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Create subnet white list", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/add/white-list",
			bytes.NewBuffer([]byte("{\"net\":\""+ipSubnet2.Subnet+"\"}")))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete subnet black list", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/delete/black-list",
			bytes.NewBuffer([]byte("{\"net\":\""+ipSubnet1.Subnet+"\"}")))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete subnet white list", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/delete/white-list",
			bytes.NewBuffer([]byte("{\"net\":\""+ipSubnet2.Subnet+"\"}")))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})

	t.Run("Delete bucket", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpReq, err := http.NewRequestWithContext(ctx, "POST", httpURL+"/delete/bucket",
			bytes.NewBuffer([]byte("{\"login\":\"test\",\"ip\":\"127.0.0.1\"}")))
		require.NoError(t, err)
		httpReq.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(httpReq)
		require.NoError(t, err)
		defer resp.Body.Close()

		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
