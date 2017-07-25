package main

import (
	"log"
	"net/http"

	"github.com/nbari/comments/github"
	"github.com/nbari/violetear"
)

var version string

func main() {
	router := violetear.New()
	router.LogRequests = true

	router.HandleFunc("/github/", github.Handler)
	router.HandleFunc("/_healthcheck_", healthCheck)

	log.Fatal(http.ListenAndServe(":8080", router))
}
