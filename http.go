package marabunta

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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

func (m *Marabunta) register(w http.ResponseWriter, r *http.Request) {
	csr, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pemBlock, _ := pem.Decode(csr)
	if pemBlock == nil {
		http.Error(w, "could not parse csr", http.StatusUnprocessableEntity)
		return
	}

	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = clientCSR.CheckSignature(); err != nil {
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	log.Printf("clientCSR = %s\n", clientCSR)

	//w.Header().Set("Content-Type", "application/x-x509-ca-cert")
	//w.Write(ca)
}
