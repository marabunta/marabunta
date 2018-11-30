package marabunta

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/gomodule/redigo/redis"
	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc"
)

// Marabunta struct
type Marabunta struct {
	clients sync.Map
	config  *Config
	db      *sql.DB
	gRPC    *grpc.Server
	redis   *redis.Pool
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

	// initialize gRPC
	gRPC, err := initGRPC(c)
	if err != nil {
		return nil, err
	}

	return &Marabunta{
		config: c,
		db:     db,
		gRPC:   gRPC,
		redis:  redis,
	}, nil
}

// Start start the services
func (m *Marabunta) Start() error {
	// listen for gRPC
	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", m.config.GRPCPort))
	if err != nil {
		return err
	}

	pb.RegisterMarabuntaServer(m.gRPC, m)

	// start gRPC server
	go func() {
		log.Fatal(m.gRPC.Serve(conn))
	}()

	// TODO gRPC events
	go m.Pulse()

	log.Printf("Starting marabunta, listening on [HTTP *:%d] [gRPC *:%d]\n", m.config.HTTPPort, m.config.GRPCPort)

	// start HTTP server
	return m.HTTP().ListenAndServeTLS(m.config.TLS.Crt, m.config.TLS.Key)
}
