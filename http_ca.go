package marabunta

import (
	"io/ioutil"
	"log"
	"net/http"
)

// HTTPCA http handler sending the CA
func (m *Marabunta) HTTPCA(w http.ResponseWriter, r *http.Request) {
	ca, err := ioutil.ReadFile(m.config.TLS.CA)
	if err != nil {
		log.Printf("HTTP CA error: %s\n", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/x-x509-ca-cert")
	w.Write(ca)
}
