package redis

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gomodule/redigo/redis"
)

// Service defines the client interface for redis usage.
type Service interface {
	// Ping runs the redis Ping command.
	Ping(ctx context.Context) (string, error)
	// Shutdown closes the redis service connection pool.
	Shutdown(ctx context.Context) error
	// Get returns the value associated with a key.
	Get(ctx context.Context, key string) (string, error)
	// Set assigns a value to a key.
	Set(ctx context.Context, key, value string) error
}

// Check in compile time that redisService implements the Service interface.
var _ Service = (*redisService)(nil)

// redisService implements the Service interface.
type redisService struct {
	pool *redis.Pool
}

// NewService is a redisService constructor.
func NewService(address string) (Service, error) {
	redisService := redisService{}

	if err := redisService.initialization(address); err != nil {
		return nil, fmt.Errorf("redis: new redis service initialization has failed: %w", err)
	}

	return &redisService, nil
}

// initialization is a helper method for service initialization.
//
// Method is used for connection pool setup with the redis server.
func (s *redisService) initialization(address string) error {
	s.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,

		DialContext: func(ctx context.Context) (redis.Conn, error) {
			c, err := redis.DialContext(ctx, "tcp", address)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	conn := s.pool.Get()
	_, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("redis: service connection has error: %w", err)
	}

	log.Println("redis: service successfully initialized, ready for use.")

	return nil
}

// Shutdown implements Service.Shutdown method.
func (s *redisService) Shutdown(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("redis: service shutdown has context error: %w", err)
	}

	if err := s.pool.Close(); err != nil {
		return fmt.Errorf("redis: service shutdown has failed: %w", err)
	}

	log.Println("redis: service shutdown is complete.")

	return nil
}

// Ping implements Servce.Ping method.
func (s *redisService) Ping(ctx context.Context) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("redis: ping command has context error: %w", err)
	}

	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return "", fmt.Errorf("redis: connection setup for ping command has failed: %w", err)
	}
	defer conn.Close()

	result, err := redis.String(conn.Do("PING"))
	if err != nil {
		return "", fmt.Errorf("redis: ping command has failed: %w", err)
	}

	return result, nil
}
