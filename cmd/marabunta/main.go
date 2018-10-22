package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/marabunta/marabunta"
	"github.com/marabunta/marabunta/pkg/github"
	"github.com/marabunta/marabunta/pkg/healthcheck"
	"github.com/nbari/violetear"
)

var version string

func main() {
	var (
		v        = flag.Bool("v", false, fmt.Sprintf("Print version: %s", version))
		httpPort = flag.Int("http", 8000, "Listen on `HTTP port`")
		grpcPort = flag.Int("grpc", 1415, "Listen on `gRPC port`")
		certFile = flag.String("cert", "server.crt", "TLS cert")
		keyFile  = flag.String("key", "server.key", "TLS key")
	)

	flag.Parse()
	if *v {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	go func() {
		err := marabunta.StartGRPC(*grpcPort, *certFile, *keyFile)
		if err != nil {
			log.Fatal(err)
		}
	}()

	// start router
	router := violetear.New()
	router.Verbose = false
	router.LogRequests = true

	router.HandleFunc("/github/", github.Handler)

	// set version on healthCheck
	healthcheck.Version = version
	router.HandleFunc("/status", healthcheck.Handler)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", *httpPort),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   7 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(srv.ListenAndServe())
}
