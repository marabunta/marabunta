package marabunta

import (
	"fmt"
	"log"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// Stream stream
func (m *Marabunta) Stream(stream pb.Marabunta_StreamServer) error {
	var client string
	ctx := stream.Context()
	if peer, ok := peer.FromContext(ctx); ok {
		tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
		client := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
		m.clients.Store(client, stream)
	}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md = %+v\n", md)
	}

	for {
		in, err := stream.Recv()
		if err != nil {
			m.clients.Delete(client)
			log.Printf("ant: %s, %s", client, err)
			return err
		}
		msg := &pb.StreamResponse{
			Event: &pb.StreamResponse_EPing{
				EPing: &pb.StreamResponse_Ping{
					Msg: fmt.Sprintf("pong: %s", in.Msg),
				},
			},
		}
		err = stream.Send(msg)
		if err != nil {
			log.Printf("ant: %s, %s", client, err)
		}
	}
}
