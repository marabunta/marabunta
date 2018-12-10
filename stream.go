package marabunta

import (
	"fmt"
	"log"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

// Stream stream
func (m *Marabunta) Stream(stream pb.Marabunta_StreamServer) error {
	var key string
	ctx := stream.Context()
	if peer, ok := peer.FromContext(ctx); ok {
		tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
		client := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
		ip := peer.Addr.String()
		key = fmt.Sprintf("%s@%s", client, ip)
		m.clients.Store(key, stream)
	}
	log.Printf("key = %+v\n", key)

	for {
		in, err := stream.Recv()
		if err != nil {
			m.clients.Delete(key)
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
			log.Printf("ant: %s, %s", key, err)
		}
	}
}
