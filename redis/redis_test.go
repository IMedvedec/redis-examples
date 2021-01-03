// +build integration

package redis_test

import (
	"context"
	"testing"

	"github.com/imedvedec/redis-examples/redis"
)

const (
	testServerAddress = ":6379"
)

// TestNewService is a connection check test for the service.
func TestNewService(t *testing.T) {
	// seccessful service setup.
	t.Run("sucessful_setup", func(t *testing.T) {
		ctx := context.Background()

		service, err := redis.NewService(testServerAddress)
		if err != nil {
			t.Fatal(err)
		}

		pong, err := service.Ping(ctx)
		if err != nil {
			t.Fatal(err)
		}

		if pong != "PONG" {
			t.Fatal("PING command result different than PONG")
		}

		err = service.Shutdown(ctx)
		if err != nil {
			t.Fatal(err)
		}
	})

	// invalid service setup.
	t.Run("invalid_setup", func(t *testing.T) {
		address := "invalidHost.com:6397"

		service, err := redis.NewService(address)

		if !(service == nil && err != nil) {
			t.Error("new service should return an initialization error with a nil poiner")
		}
	})
}
