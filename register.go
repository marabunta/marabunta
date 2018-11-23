package marabunta

import (
	"context"
	"errors"

	pb "github.com/marabunta/protobuf"
)

// Register register
func (m *Marabunta) Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return nil, errors.New("could not find commonName from TLS")
}
