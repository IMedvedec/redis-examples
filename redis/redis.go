package redis

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	redisHostEnvVarKey = "REDIS_HOST"
	defaultRedisHost   = ":6379"
)

// Service defines the client interface for redis usage.
type Service interface {
	// Ping runs the redis Ping command.
	Ping(ctx context.Context) (string, error)
	// Shutdown closes the redis service connection pool.
	Shutdown(ctx context.Context) error
}

// Check in compile time that redisService implements the Service interface.
var _ Service = (*redisService)(nil)

// redisService implements the Service interface.
type redisService struct {
	pool *redis.Pool
}

// NewService is a redisService constructor.
func NewService() (Service, error) {
	redisService := redisService{}

	if err := redisService.initialization(); err != nil {
		return nil, fmt.Errorf("redis: new redis service initialization has failed: %w", err)
	}

	return &redisService, nil
}

// initialization is a helper method for service initialization.
//
// Method is used for connection pool setup with the redis server.
func (s *redisService) initialization() error {
	redisHost := os.Getenv(redisHostEnvVarKey)
	if redisHost == "" {
		redisHost = defaultRedisHost
	}

	s.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,

		DialContext: func(ctx context.Context) (redis.Conn, error) {
			c, err := redis.DialContext(ctx, "tcp", redisHost)
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
		return "", fmt.Errorf("redis: connection setup for PING command has failed: %w", err)
	}
	defer conn.Close()

	result, err := redis.String(conn.Do("PING"))
	if err != nil {
		return "", fmt.Errorf("redis: ping command has failed: %w", err)
	}

	return result, nil
}
