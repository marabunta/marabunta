package marabunta

import (
	"errors"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func (m *Marabunta) streamInterceptor(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	ctx := stream.Context()
	if peer, ok := peer.FromContext(ctx); ok {
		tlsInfo := peer.AuthInfo.(credentials.TLSInfo)
		client := tlsInfo.State.VerifiedChains[0][0].Subject.CommonName
		ip := peer.Addr.String()

		// TODO update ant's table
		key := fmt.Sprintf("%s@%s", client, ip)
		fmt.Printf("client = %s\n", client)
		fmt.Printf("ip.Addr = %s\n", ip)
		fmt.Printf("key for map = %+v\n", key)
		// read metadata
		if md, ok := metadata.FromIncomingContext(ctx); ok {
			log.Printf("md = %+v\n", md)
		}

		return handler(srv, stream)
	}
	return errors.New("could not get client info")
}
