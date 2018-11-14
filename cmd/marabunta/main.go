package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marabunta/marabunta"
)

var version string

func main() {
	parser := &marabunta.Parse{}

	fs := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	fs.Usage = parser.Usage(fs)

	cfg, err := parser.ParseArgs(fs)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if (fs.Lookup("v")).Value.(flag.Getter).Get().(bool) {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}
	fmt.Printf("cfg = %+v\n", cfg)

	m, err := marabunta.New(cfg)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	m.StartGRPC()
	m.StartHTTP()

	//go func() {
	//err := marabunta.StartGRPC(*grpcPort, *certFile, *keyFile)
	//if err != nil {
	//log.Fatal(err)
	//}
	//}()

	//// start router
	//router := violetear.New()
	//router.Verbose = false
	//router.LogRequests = true

	//router.HandleFunc("/github/", github.Handler)

	//// set version on healthCheck
	//healthcheck.Version = version
	//router.HandleFunc("/status", healthcheck.Handler)

	//srv := &http.Server{
	//Addr:           fmt.Sprintf(":%d", *httpPort),
	//Handler:        router,
	//ReadTimeout:    5 * time.Second,
	//WriteTimeout:   7 * time.Second,
	//MaxHeaderBytes: 1 << 20,
	//}
	//log.Fatal(srv.ListenAndServe())
}
