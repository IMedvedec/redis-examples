package redis

import (
	"context"
	"errors"
	"fmt"

	"github.com/gomodule/redigo/redis"
)

var (
	NotExistErr = errors.New("string value does not exist")
)

// Get implements Service.Get method.
func (s *redisService) Get(ctx context.Context, key string) (string, error) {
	if err := ctx.Err(); err != nil {
		return "", fmt.Errorf("redis: get command has context error: %w", err)
	}

	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return "", fmt.Errorf("redis: connection setup for get command has failed: %w", err)
	}
	defer conn.Close()

	result, err := redis.String(conn.Do("GET", key))
	if err != nil {
		if errors.Is(err, redis.ErrNil) {
			return "", fmt.Errorf("redis: string value does not exist for key: %w", NotExistErr)
		}
		return "", fmt.Errorf("redis: get command has failed: %w", err)
	}

	return result, nil
}

// Set implements Service.Set method.
func (s *redisService) Set(ctx context.Context, key, value string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("redis: set command has context error: %w", err)
	}

	conn, err := s.pool.GetContext(ctx)
	if err != nil {
		return fmt.Errorf("redis: connection setup for set command has failed: %w", err)
	}
	defer conn.Close()

	_, err = redis.String(conn.Do("SET", key, value))
	if err != nil {
		return fmt.Errorf("redis: set command has failed: %w", err)
	}

	return nil
}
