package encryption

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

type ECDSAHelper struct{}

func (h ECDSAHelper) Encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func (h ECDSAHelper) Decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		println(err.Error())
	}

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))

	x509EncodedPub := blockPub.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509EncodedPub)
	if err != nil {
		println(err.Error())
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}

func (h ECDSAHelper) DecodePrivateKeyFromPem(pemEncoded string) *ecdsa.PrivateKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from pem: %v", err.Error())
		return nil
	}

	return privateKey
}

func (h ECDSAHelper) DecodePublicKeyFromPem(pemEncoded string) *ecdsa.PublicKey {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	genericPublicKey, err := x509.ParsePKIXPublicKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from pem: %v", err.Error())
		return nil
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey
}

func (h ECDSAHelper) DecodePrivateKeyFromBase64(bas64Encoded string) *ecdsa.PrivateKey {
	x509Encoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}

	privateKey, err := x509.ParseECPrivateKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	return privateKey
}

func (h ECDSAHelper) DecodePublicKeyFromBase64(bas64Encoded string) *ecdsa.PublicKey {
	x509Encoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}

	genericPublicKey, err := x509.ParsePKIXPublicKey(x509Encoded)
	if err != nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	return publicKey
}

func (h ECDSAHelper) DecodePrivateKeyFromBase64Pem(bas64Encoded string) *ecdsa.PrivateKey {
	pemEncoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}
	privateKey := h.DecodePrivateKeyFromPem(string(pemEncoded))

	// privateKey, err := x509.ParseECPrivateKey(pemEncoded)
	if privateKey == nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	return privateKey
}

func (h ECDSAHelper) DecodePublicKeyFromBase64Pem(bas64Encoded string) *ecdsa.PublicKey {
	pemEncoded, err := base64.URLEncoding.DecodeString(bas64Encoded)
	if err != nil {
		logger.Error("There was an error decoding the private key from base64: %v", err.Error())
		return nil
	}
	privateKey := h.DecodePublicKeyFromPem(string(pemEncoded))

	// privateKey, err := x509.ParseECPrivateKey(pemEncoded)
	if privateKey == nil {
		logger.Error("There was an error getting the private key from base64: %v", err.Error())
		return nil
	}

	return privateKey
}

func (h ECDSAHelper) GeneratePrivateKey(size EncryptionKeySize) *ecdsa.PrivateKey {
	if size > 512 {
		size = Bit512
	}

	var priv *ecdsa.PrivateKey
	var err error
	switch size {
	case Bit256:
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case Bit384:
		priv, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case Bit512:
		priv, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	default:
		priv, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	}

	if err != nil {
		return nil
	}

	return priv
}

func (h ECDSAHelper) GenerateKeys(size EncryptionKeySize) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	priv := h.GeneratePrivateKey(size)

	if priv == nil {
		return nil, nil
	}

	return priv, &priv.PublicKey
}
