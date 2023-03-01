package encryption

import (
	"crypto/ecdsa"
	"crypto/md5"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"

	log "github.com/cjlapao/common-go-logger"
)

var logger = log.Get()

func GetKeyBytes(key interface{}) []byte {
	var encodedBytes []byte
	switch k := key.(type) {
	case *rsa.PrivateKey:
		encodedBytes = x509.MarshalPKCS1PrivateKey(k)
		return encodedBytes
	case *rsa.PublicKey:
		encodedBytes = x509.MarshalPKCS1PublicKey(k)
		return encodedBytes
	case *ecdsa.PrivateKey:
		encodedBytes, _ = x509.MarshalECPrivateKey(k)
		return encodedBytes
	case *ecdsa.PublicKey:
		encodedBytes, _ = x509.MarshalPKIXPublicKey(k)
		return encodedBytes
	default:
		return make([]byte, 0)
	}
}

func GetSHA2KeyFingerprint(key interface{}) [32]byte {
	encodedBytes := GetKeyBytes(key)
	fingerprint := sha256.Sum256(encodedBytes)
	return fingerprint
}

func GetSHA1KeyFingerprint(key interface{}) [20]byte {
	encodedBytes := GetKeyBytes(key)
	fingerprint := sha1.Sum(encodedBytes)
	return fingerprint
}

func GetMD5KeyFingerprint(key interface{}) [16]byte {
	encodedBytes := GetKeyBytes(key)
	fingerprint := md5.Sum(encodedBytes)
	return fingerprint
}

func GetBase64KeyFingerprint(key interface{}) string {
	fingerprintBytes := GetSHA2KeyFingerprint(key)
	encoded := base64.URLEncoding.EncodeToString(fingerprintBytes[:])
	return encoded
}
