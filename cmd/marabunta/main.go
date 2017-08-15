package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/marabunta/marabunta/pkg/github"
	"github.com/marabunta/marabunta/pkg/healthcheck"
	"github.com/nbari/violetear"
)

var version string

func main() {
	v := flag.Bool("v", false, fmt.Sprintf("Print version: %s", version))

	flag.Parse()
	if *v {
		fmt.Printf("%s\n", version)
		os.Exit(0)
	}

	// start router
	router := violetear.New()
	router.LogRequests = true

	router.HandleFunc("/github/", github.Handler)

	// set version on healthCheck
	healthcheck.Version = version
	router.HandleFunc("/_healthcheck_", healthcheck.Handler)

	log.Fatal(http.ListenAndServe(":8080", router))
}
