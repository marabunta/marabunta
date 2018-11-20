package marabunta

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net"
	"sync"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Server represents the gRPC server
type Server struct {
	marabunta sync.Map
}

func (s *Server) Update(ctx context.Context, update *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return &pb.UpdateResponse{Ok: true}, nil
}

// StartGRPC start gRPC server
func StartGRPC(port int, ca, crt, key string) error {
	cert, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return err
	}

	CA, err := ioutil.ReadFile(ca)
	if err != nil {
		return fmt.Errorf("could not read CA certificate: %s", err)
	}

	// Append the client certificates from the CA
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(CA); !ok {
		return fmt.Errorf("failed to append client certs")
	}

	tlsConfig := &tls.Config{
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	)

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
