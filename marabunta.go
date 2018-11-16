package marabunta

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"sync"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type marabunta struct {
	db    *sql.DB
	redis string
}

// Server represents the gRPC server
type Server struct {
	marabunta sync.Map
}

func (s *Server) Update(ctx context.Context, update *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return &pb.UpdateResponse{Ok: true}, nil
}

// StartGRPC start gRPC server
func StartGRPC(port int, cert, key string) error {
	creds, err := credentials.NewServerTLSFromFile("server.crt", "server.key")
	if err != nil {
		log.Fatalf("could not load TLS keys: %s", err)
	}

	// Create an array of gRPC options with the credentials
	opts := []grpc.ServerOption{grpc.Creds(creds)}

	grpcServer := grpc.NewServer(opts...)
	marabunta := &Server{}
	pb.RegisterMarabuntaServer(grpcServer, marabunta)

	conn, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	// TODO events
	go marabunta.Pulse()

	return grpcServer.Serve(conn)
}

func New(c *Config) (*marabunta, error) {
	db, err := initMySQL(c)
	if err != nil {
		return nil, err
	}

	return &marabunta{
		db: db,
	}, nil
}

func (m *marabunta) Start() error {
	return nil
}
