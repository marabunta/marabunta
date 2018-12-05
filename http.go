package marabunta

import (
	"fmt"
	"net/http"
	"time"

	"github.com/marabunta/marabunta/http/healthcheck"
	"github.com/nbari/violetear"
)

// HTTP returns http router
func (m *Marabunta) HTTP() *http.Server {
	// HTTP router
	router := violetear.New()
	router.Verbose = false
	router.LogRequests = true

	router.HandleFunc("/register", m.register, "POST")

	// set version on healthCheck
	healthcheck.Version = "foo"
	router.HandleFunc("/status", healthcheck.Handler)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", m.config.HTTPPort),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   7 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return srv
}
