package marabunta

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"sync"

	"github.com/gomodule/redigo/redis"
	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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

	return &Marabunta{
		config: c,
		db:     db,
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

	cert, err := tls.LoadX509KeyPair(m.config.TLS.Crt, m.config.TLS.Key)
	if err != nil {
		return err
	}

	caCert, err := ioutil.ReadFile(m.config.TLS.CACrt)
	if err != nil {
		return fmt.Errorf("could not read CA certificate from file %q, %s", m.config.TLS.CACrt, err)
	}

	// Append the client certificates from the CA
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return errors.New("could not append CA certificate to cert pool")
	}

	// mTLS
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	// create a gRPC server and get client info from the interceptors
	m.gRPC = grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
		grpc.StreamInterceptor(m.streamInterceptor),
	)

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
