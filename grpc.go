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

	caCert, err := ioutil.ReadFile(c.TLS.CA)
	if err != nil {
		return nil, fmt.Errorf("could not read CA certificate from file %q, %s", c.TLS.CA, err)
	}

	// Append the client certificates from the CA
	caCertPool := x509.NewCertPool()
	if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
		return nil, errors.New("could not append CA certificate to cert pool")
	}

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    caCertPool,
	}

	return grpc.NewServer(
		grpc.Creds(credentials.NewTLS(tlsConfig)),
	), nil
}
