package marabunta

import (
	"context"

	pb "github.com/marabunta/protobuf"
)

// Payload grpc
func (m *Marabunta) Payload(ctx context.Context, update *pb.PayloadRequest) (*pb.PayloadResponse, error) {
	return nil, nil
}
