package marabunta

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/marabunta/marabunta/pkg/github"
	"github.com/marabunta/marabunta/pkg/healthcheck"
	pb "github.com/marabunta/protobuf"
	"github.com/nbari/violetear"
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

	// TODO events
	go m.Pulse()

	// start gRPC server
	go func() {
		log.Fatal(m.gRPC.Serve(conn))
	}()

	// HTTP router
	router := violetear.New()
	router.Verbose = false
	router.LogRequests = true

	router.HandleFunc("/github/", github.Handler)

	// set version on healthCheck
	healthcheck.Version = "foo"
	router.HandleFunc("/status", healthcheck.Handler)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", m.config.HTTPPort),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   7 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// start HTTP server
	return srv.ListenAndServe()
}
