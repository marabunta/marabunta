package marabunta

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func initRedis(c *Config) (*redis.Pool, error) {
	pool := &redis.Pool{
		MaxIdle:     10,
		IdleTimeout: 60 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%s:%d",
				c.Redis.Host,
				c.Redis.Port))
		},
	}

	// Get a connection
	conn := pool.Get()
	defer conn.Close()

	// Test the connection
	_, err := conn.Do("PING")
	if err != nil {
		return nil, fmt.Errorf("Can't connect to the Redis database: %s", err)
	}
	return pool, nil
}
