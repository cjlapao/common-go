package authorization_context

import (
	"time"

	"github.com/cjlapao/common-go/security/encryption"
)

type AuthorizationOptions struct {
	Issuer        string
	Audiences     []string
	Scope         string
	TokenDuration time.Duration
	SignatureType encryption.EncryptionKeyType
	SignatureSize encryption.EncryptionKeySize
	PrivateKey    string
	PublicKey     string
	KeyId         string
}

type AuthorizationValidationOptions struct {
	Audiences     bool
	ExpiryDate    bool
	Subject       bool
	Issuer        bool
	VerifiedEmail bool
	NotBefore     bool
	Tenant        bool
}
