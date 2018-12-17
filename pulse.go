package marabunta

import (
	"fmt"
	"log"
	"sync"
	"time"

	pb "github.com/marabunta/protobuf"
)

func (m *Marabunta) Pulse() {
	getClients := func(ant, streams interface{}) bool {
		send := func(ant, stream interface{}) bool {
			log.Printf("ant = %+v\n", ant)
			msg := &pb.StreamResponse{
				Event: &pb.StreamResponse_EPulse{
					EPulse: &pb.StreamResponse_Pulse{
						Msg: fmt.Sprintf("to: %s msg: %s", ant, "get pulse every 10 seconds?"),
					},
				},
			}
			err := stream.(pb.Marabunta_StreamServer).Send(msg)
			if err != nil {
				log.Printf("ant: %s, %s", ant, err)
				return false
			}
			return true
		}
		streams.(*sync.Map).Range(send)
		return true
	}
	for {
		select {
		case <-time.After(10 * time.Second):
			m.clients.Range(getClients)
		}
	}
}
