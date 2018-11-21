package marabunta

import (
	"context"

	pb "github.com/marabunta/protobuf"
)

// Update grpc
func (m *Marabunta) Update(ctx context.Context, update *pb.UpdateRequest) (*pb.UpdateResponse, error) {
	return &pb.UpdateResponse{Ok: true}, nil
}
