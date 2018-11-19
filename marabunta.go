package marabunta

import (
	"database/sql"

	"github.com/gomodule/redigo/redis"
)

// Marabunta struct
type Marabunta struct {
	db     *sql.DB
	redis  *redis.Pool
	config *Config
	errc   chan error
}

// New return a marabunta
func New(c *Config) (*Marabunta, error) {
	// initialize MySQL create databases if needed
	db, err := initMySQL(c)
	if err != nil {
		return nil, err
	}

	// initialize Redis
	redis, err := initRedis(c)
	if err != nil {
		return nil, err
	}

	return &Marabunta{
		db:     db,
		redis:  redis,
		config: c,
		errc:   make(chan error),
	}, nil
}

// Start start the services
func (m *Marabunta) Start() error {
	go func() {
		err := StartGRPC(m.config.GRPCPort, m.config.TLS.Crt, m.config.TLS.Key)
		if err != nil {
			m.errc <- err
		}
	}()

	// TODO
	select {
	case err := <-m.errc:
		return err
	}
}
