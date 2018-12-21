package marabunta

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"io"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"
)

func (m *Marabunta) register(w http.ResponseWriter, r *http.Request) {
	csr, err := ioutil.ReadAll(io.LimitReader(r.Body, 4096))
	if err != nil {
		// 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pemBlock, _ := pem.Decode(csr)
	if pemBlock == nil {
		// 422
		http.Error(w, "could not parse csr", http.StatusUnprocessableEntity)
		return
	}

	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		// 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err = clientCSR.CheckSignature(); err != nil {
		// 406
		http.Error(w, err.Error(), http.StatusNotAcceptable)
		return
	}

	caFile, err := ioutil.ReadFile(m.config.TLS.CACrt)
	if err != nil {
		// 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pemBlock, _ = pem.Decode(caFile)
	if pemBlock == nil {
		// 500
		http.Error(w, "could not parse csr", http.StatusInternalServerError)
		return
	}

	caCRT, err := x509.ParseCertificate(pemBlock.Bytes)
	if err != nil {
		// 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	caPrivateKeyFile, err := ioutil.ReadFile(m.config.TLS.CAKey)
	if err != nil {
		// 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pemBlock, _ = pem.Decode(caPrivateKeyFile)
	if pemBlock == nil {
		// 500
		http.Error(w, "could not parse CA.key", http.StatusInternalServerError)
		return
	}

	caPrivateKey, err := x509.ParseECPrivateKey(pemBlock.Bytes)
	if err != nil {
		// 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// create client certificate template
	clientCRTTemplate := x509.Certificate{
		ExtKeyUsage:        []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		Issuer:             caCRT.Subject,
		KeyUsage:           x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		NotAfter:           time.Now().AddDate(3, 0, 0),
		NotBefore:          time.Now(),
		PublicKey:          clientCSR.PublicKey,
		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		SerialNumber:       big.NewInt(2),
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,
		Subject:            clientCSR.Subject,
	}

	// create client certificate from template and CA public key
	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, caCRT, clientCSR.PublicKey, caPrivateKey)
	if err != nil {
		// 500
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	crt := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: clientCRTRaw})

	w.Header().Set("Content-Type", "application/x-x509-ca-cert")
	w.Write(append(crt, caFile...))
}
