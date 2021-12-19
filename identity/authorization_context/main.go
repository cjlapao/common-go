package authorization_context

import (
	"errors"
	"strings"
	"time"

	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/security"
	"github.com/google/uuid"
)

type AuthorizationContext struct {
	User          *ContextUser
	CorrelationId string
	Options       AuthorizationOptions
}

var currentAuthorizationContext *AuthorizationContext

func NewFromUser(user *ContextUser) *AuthorizationContext {
	currentAuthorizationContext = &AuthorizationContext{
		User:          user,
		CorrelationId: uuid.NewString(),
	}

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
	publicKey := config.GetString("JWT" + keyId + "_PUBLIC_KEY")

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

	if privateKey == "" {
		panic(errors.New("private key not found"))
	}
	switch strings.ToLower(authorizationType) {
	case "hmac":
		if privateKey == "" {
			privateKey = "SomeRandomSecret"
		}
		if keySize == "" {
			keySize = "256bit"
		}
		a.Options.SignatureType = security.HMAC
		a.Options.SignatureSize = a.Options.SignatureSize.FromString(keySize)
		a.Options.PrivateKey = tryDecodekey(privateKey)
	case "ecdsa":
		if publicKey == "" {
			panic(errors.New("public key not found"))
		}
		if keySize == "" {
			keySize = "256bit"
		}
		a.Options.SignatureType = security.ECDSA
		a.Options.SignatureSize = a.Options.SignatureSize.FromString(keySize)
		a.Options.PrivateKey = tryDecodekey(privateKey)
		a.Options.PublicKey = tryDecodekey(publicKey)
	case "rsa":
		if publicKey == "" {
			panic(errors.New("public key not found"))
		}
		if keySize == "" {
			keySize = "1024bit"
		}
		a.Options.SignatureType = security.RSA
		a.Options.SignatureSize = a.Options.SignatureSize.FromString(keySize)
		a.Options.PrivateKey = tryDecodekey(privateKey)
		a.Options.PublicKey = tryDecodekey(publicKey)
	}

	return a
}

func GetCurrent() *AuthorizationContext {
	if currentAuthorizationContext != nil {
		return currentAuthorizationContext
	}

	return nil
}

func tryDecodekey(value string) string {
	decoded, err := security.DecodeBase64String(value)
	if err == nil {
		return decoded
	} else {
		return value
	}
}
