package marabunta

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func initGRPC(c *Config) (*grpc.Server, error) {
	cert, err := tls.LoadX509KeyPair(c.TLS.Crt, c.TLS.Key)
	if err != nil {
		return nil, err
	}

	caCert, err := ioutil.ReadFile(c.TLS.CACrt)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate from file %q, %s", c.TLS.CACrt, err)
	}

	// Append the client certificates from the CA
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not append CA certificate to cert pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	}

	return grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	), nil
}
