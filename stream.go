package marabunta

import (
	"fmt"
	"log"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc/metadata"
)

// Stream stream
func (s *Server) Stream(stream pb.Marabunta_StreamServer) error {
	var ant string

	if md, ok := metadata.FromIncomingContext(stream.Context()); ok {
		log.Printf("md = %+v\n", md)
		ant = md["ant"][0]
		s.marabunta.Store(ant, stream)
	}

	for {
		in, err := stream.Recv()
		if err != nil {
			s.marabunta.Delete(ant)
			log.Printf("ant: %s, %s", ant, err)
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
			log.Printf("ant: %s, %s", ant, err)
		}
	}
}
