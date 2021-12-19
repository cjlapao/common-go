package authorization_context

import (
	"time"

	"github.com/cjlapao/common-go/security"
)

type AuthorizationOptions struct {
	Issuer        string
	Scope         string
	TokenDuration time.Duration
	SignatureType security.EncryptionKeyType
	SignatureSize security.EncryptionKeySize
	PrivateKey    string
	PublicKey     string
	KeyId         string
}
