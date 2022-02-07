package authorization_context

import (
	"github.com/cjlapao/common-go/security/encryption"
)

type AuthorizationOptions struct {
	KeyVaultEnabled          bool
	TokenDuration            int
	RefreshTokenDuration     int
	VerifyEmailTokenDuration int
	SignatureType            encryption.EncryptionKeyType
	SignatureSize            encryption.EncryptionKeySize
	PrivateKey               string
	PublicKey                string
	KeyId                    string
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
