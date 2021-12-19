package security

import "strings"

type EncryptionKeyType int64

const (
	ECDSA EncryptionKeyType = iota
	HMAC
	RSA
)

func (j EncryptionKeyType) String() string {
	switch j {
	case ECDSA:
		return "ecdsa"
	case HMAC:
		return "hmac"
	case RSA:
		return "rsa"
	default:
		return "hmac"
	}
}

func (j EncryptionKeyType) FromString(keyType string) EncryptionKeyType {
	switch strings.ToLower(keyType) {
	case "ecdsa":
		return ECDSA
	case "hmac":
		return HMAC
	case "rsa":
		return RSA
	default:
		return HMAC
	}
}
