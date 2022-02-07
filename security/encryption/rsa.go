package encryption

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type RSAHelper struct{}

func (h RSAHelper) Encode(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) (string, string) {
	x509Encoded := x509.MarshalPKCS1PrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub := x509.MarshalPKCS1PublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func (h RSAHelper) Decode(pemEncoded string, pemEncodedPub string) (*rsa.PrivateKey, *rsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParsePKCS1PrivateKey(x509Encoded)
	if err != nil {
		println(err.Error())
	}

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))

	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		println(err.Error())
	}

	publicKey := genericPublicKey.(*rsa.PublicKey)

	return privateKey, publicKey
}

func (h RSAHelper) DecodePrivateKeyFromPem(pemEncoded string) *rsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParsePKCS1PrivateKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from pem: %v", err.Error())
		return nil
	}

	return privateKey
}

func (h RSAHelper) DecodePublicKeyFromPem(pemEncoded string) *rsa.PublicKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	publicKey, err := x509.ParsePKCS1PublicKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from pem: %v", err.Error())
		return nil
	}

	return publicKey
}

func (h RSAHelper) DecodePrivateKeyFromBase64(bas64Encoded string) *rsa.PrivateKey {
	x509Encoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}

	privateKey, err := x509.ParsePKCS1PrivateKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	return privateKey
}

func (h RSAHelper) DecodePublicKeyFromBase64(bas64Encoded string) *rsa.PublicKey {
	x509Encoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}

	publicKey, err := x509.ParsePKCS1PublicKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	return publicKey
}

func (h RSAHelper) DecodePrivateKeyFromBase64Pem(bas64Encoded string) *rsa.PrivateKey {
	pemEncoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}

	privateKey := h.DecodePrivateKeyFromPem(string(pemEncoded))
	if privateKey == nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	return privateKey
}

func (h RSAHelper) DecodePublicKeyFromBase64Pem(bas64Encoded string) *rsa.PublicKey {
	pemEncoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}

	publicKey := h.DecodePublicKeyFromPem(string(pemEncoded))
	if publicKey == nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	return publicKey
}

func (h RSAHelper) GeneratePrivateKey(size EncryptionKeySize) *rsa.PrivateKey {
	priv, err := rsa.GenerateKey(rand.Reader, int(size))

	if err != nil {
		return nil
	}

	return priv
}

func (h RSAHelper) GenerateKeys(size EncryptionKeySize) (*rsa.PrivateKey, *rsa.PublicKey) {
	priv := h.GeneratePrivateKey(size)

	if priv == nil {
		return nil, nil
	}

	return priv, &priv.PublicKey
}
