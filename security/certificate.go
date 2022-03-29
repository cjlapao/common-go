package security

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"math/big"
	"time"
)

type pkcs1PublicKey struct {
	N *big.Int
	E int
}

type SelfSignedCertificate struct {
}

func NewRsaSelfSignedCertificate(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) *x509.Certificate {
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 180),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	template2 := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Acme C1o"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(time.Hour * 24 * 260),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey, privateKey)
	derBytes1, err := x509.CreateCertificate(rand.Reader, &template2, &template2, publicKey, privateKey)

	if err != nil {
		logger.Error("There was an error generating the self signed certificate: %v", err.Error())
	}

	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	certTemplate1Fingerprint := sha256.Sum256(derBytes)
	certTemplate2Fingerprint := sha256.Sum256(derBytes1)
	str1 := base64.StdEncoding.EncodeToString(certTemplate1Fingerprint[:])
	str2 := base64.StdEncoding.EncodeToString(certTemplate2Fingerprint[:])
	println(str1)
	println(str2)

	if err != nil {
		logger.Error("There was an error generating the self signed certificate: %v", err.Error())
	}
	fmt.Println(out.String())
	cert, err := x509.ParseCertificate(derBytes)
	return cert
}
