package marabunta

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/marabunta/marabunta/http/healthcheck"
	"github.com/marabunta/marabunta/http/register"
	"github.com/nbari/violetear"
)

// HTTP returns http router
func (m *Marabunta) HTTP() *http.Server {
	// HTTP router
	router := violetear.New()
	router.Verbose = false
	router.LogRequests = true

	router.HandleFunc("/register/", register.POST, "POST")

	// set version on healthCheck
	healthcheck.Version = "foo"
	router.HandleFunc("/status", healthcheck.Handler)
	router.HandleFunc("/ca", m.httpCA, "GET")

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", m.config.HTTPPort),
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   7 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	return srv
}

func (m *Marabunta) httpCA(w http.ResponseWriter, r *http.Request) {
	ca, err := ioutil.ReadFile(m.config.TLS.CA)
	if err != nil {
		log.Printf("HTTP CA error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")
	w.Write(ca)
}
