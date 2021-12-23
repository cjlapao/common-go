package authorization_context

import (
	"errors"
	"strings"
	"time"

	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/identity/jwt_keyvault"
	"github.com/cjlapao/common-go/security/encryption"
	"github.com/google/uuid"
)

type AuthorizationContext struct {
	User              *ContextUser
	CorrelationId     string
	TenantId          string
	Options           AuthorizationOptions
	ValidationOptions AuthorizationValidationOptions
	KeyVault          *jwt_keyvault.JwtKeyVaultService
}

var currentAuthorizationContext *AuthorizationContext

func NewFromUser(user *ContextUser) *AuthorizationContext {
	currentAuthorizationContext = &AuthorizationContext{
		User:          user,
		CorrelationId: uuid.NewString(),
		ValidationOptions: AuthorizationValidationOptions{
			Audiences:     false,
			ExpiryDate:    true,
			Subject:       true,
			Issuer:        true,
			VerifiedEmail: false,
			NotBefore:     false,
		},
	}

	currentAuthorizationContext.KeyVault = jwt_keyvault.Get()
	currentAuthorizationContext.WithDefaultOptions()

	return currentAuthorizationContext
}

func New() *AuthorizationContext {
	user := NewContextUser()

	return NewFromUser(user)
}

func (a *AuthorizationContext) WithOptions(options AuthorizationOptions) *AuthorizationContext {
	a.Options = options
	return a
}

func (a *AuthorizationContext) WithDefaultOptions() *AuthorizationContext {
	config := configuration.Get()
	issuer := config.GetString("JWT_ISSUER")
	tokenDuration := config.GetInt("JWT_TOKEN_DURATION")
	tokenScope := config.GetString("JWT_SCOPE")
	authorizationType := config.GetString("JWT_AUTH_TYPE")
	keySize := config.GetString("JWT_KEY_SIZE")
	keyId := config.GetString("JWT_KEY_ID")
	if keyId != "" {
		keyId = "_" + keyId
	}
	privateKey := config.GetString("JWT" + keyId + "_PRIVATE_KEY")

	if issuer == "" {
		issuer = "localhost"
	}

	if tokenDuration == 0 {
		tokenDuration = 60
	}

	if tokenScope == "" {
		tokenScope = "authorization"
	}

	if authorizationType == "" {
		authorizationType = "hmac"
	}

	a.Options = AuthorizationOptions{
		Issuer:        issuer,
		TokenDuration: time.Minute * time.Duration(tokenDuration),
		Scope:         tokenScope,
	}

	a.Options.Audiences = make([]string, 0)
	a.Options.Audiences = append(a.Options.Audiences, "http://localhost")

	if privateKey == "" {
		panic(errors.New("private key not found"))
	}
	keyId = strings.TrimLeft(keyId, "_")

	switch strings.ToLower(authorizationType) {
	case "hmac":
		if keySize == "" {
			keySize = "256bit"
		}
		var size encryption.EncryptionKeySize
		size = size.FromString(keySize)
		a.KeyVault.WithBase64HmacKey(keyId, privateKey, size)
	case "ecdsa":
		a.KeyVault.WithBase64EcdsaKey(keyId, privateKey)
	case "rsa":
		a.KeyVault.WithBase64RsaKey(keyId, privateKey)
	}

	return a
}

func (a *AuthorizationContext) WithAudience(audience string) *AuthorizationContext {
	found := false
	for _, inAudience := range a.Options.Audiences {
		if strings.EqualFold(inAudience, audience) {
			found = true
			break
		}
	}
	if !found {
		a.Options.Audiences = append(a.Options.Audiences, audience)
	}

	return a
}

func GetCurrent() *AuthorizationContext {
	if currentAuthorizationContext != nil {
		return currentAuthorizationContext
	}

	return nil
}
