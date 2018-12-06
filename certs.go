package marabunta

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io/ioutil"
	"math/big"
	"path/filepath"
	"time"
)

func createCertificates(cfg *Config) error {
	// create client certificate template
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().Unix()),
		Subject: pkix.Name{
			Organization: []string{"marabunta"},
			CommonName:   "marabunta",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}

	pub := &priv.PublicKey

	caCrt, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type: "CERTIFICATE",
		Headers: map[string]string{
			"CA": "marabunta",
		},
		Bytes: caCrt,
	}
	err = ioutil.WriteFile(filepath.Join(cfg.Home, "CA.crt"), pem.EncodeToMemory(block), 0644)
	if err != nil {
		return err
	}

	privKey, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return err
	}

	block = &pem.Block{Type: "PRIVATE KEY", Bytes: privKey}
	err = ioutil.WriteFile(filepath.Join(cfg.Home, "CA.key"), pem.EncodeToMemory(block), 0644)
	if err != nil {
		return err
	}

	return nil
}
