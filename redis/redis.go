package redis

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

const (
	redisHostEnvVarKey = "REDIS_HOST"
)

// Service defines the client interface for redis usage.
type Service interface {
	// Ping runs the redis Ping command.
	Ping() error
	// Shutdown closes the redis service connection pool.
	Shutdown() error
}

// Check in compile time that redisService implements the Service interface.
var _ Service = (*redisService)(nil)

// redisService implements the Service interface.
type redisService struct {
	pool *redis.Pool
}

// NewService is a redisService constructor.
func NewService() Service {
	redisService := redisService{}

	if err := redisService.initialization(); err != nil {
		log.Panicf("redis: new redis service initialization has failed: %v", err)
	}

	return &redisService
}

// initialization is a helper method for service initialization.
//
// Method is used for connection pool setup with the redis server.
func (s *redisService) initialization() error {
	redisHost := os.Getenv(redisHostEnvVarKey)
	if redisHost == "" {
		redisHost = ":6379"
	}

	s.pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 180 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", redisHost)
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

	log.Println("redis: service successfully initialized, ready for use.")

	return nil
}

// Shutdown implements Service.Shutdown method.
func (s *redisService) Shutdown() error {
	s.pool.Close()

	log.Println("redis: service shutdown is complete.")

	return nil
}

// Ping implements Servce.Ping method.
func (s *redisService) Ping() error {
	conn := s.pool.Get()
	defer conn.Close()

	res, err := redis.String(conn.Do("PING"))
	if err != nil {
		return fmt.Errorf("redis: ping command has error: %w", err)
	}
	log.Println(res)

	return nil
}
