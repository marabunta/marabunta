package marabunta

import (
	"errors"
	"fmt"
	"log"
	"sync"

	pb "github.com/marabunta/protobuf"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

func (m *Marabunta) addStream(client, ip string, stream pb.Marabunta_StreamServer) {
	ant, ok := m.clients.Load(client)
	if ok {
		ant.(*sync.Map).Store(ip, stream)
	} else {
		conns := &sync.Map{}
		conns.Store(ip, stream)
		m.clients.Store(client, conns)
	}
}

func (m *Marabunta) delStream(client, ip string) {
	ant, ok := m.clients.Load(client)
	if ok {
		ant.(*sync.Map).Delete(ip)
		length := 0
		ant.(*sync.Map).Range(func(_, _ interface{}) bool {
			length++
			return true
		})
		if length == 0 {
			log.Printf("removing client: %s\n", client)
			m.clients.Delete(client)
		}
	}
}

// Stream stream
func (m *Marabunta) Stream(stream pb.Marabunta_StreamServer) error {
	var (
		client string
		ip     string
	)
	ctx := stream.Context()
	if peer, ok := peer.FromContext(ctx); ok {
		tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
		client = tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
		ip = peer.Addr.String()
		m.addStream(client, ip, stream)
	} else {
		return errors.New("could not read peer from Context")
	}

	for {
		in, err := stream.Recv()
		if err != nil {
			log.Printf("removing [%s %q], %s", client, ip, err)
			m.delStream(client, ip)
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
			log.Printf("error sending to [%s %q], %s", client, ip, err)
		}
	}
}
