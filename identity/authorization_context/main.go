package authorization_context

import (
	"errors"
	"net/http"
	"strings"

	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/identity/interfaces"
	"github.com/cjlapao/common-go/identity/jwt_keyvault"
	"github.com/cjlapao/common-go/security/encryption"
	"github.com/cjlapao/common-go/service_provider"
)

type AuthorizationContext struct {
	User              *ContextUser
	TenantId          string
	Issuer            string
	Scope             string
	Audiences         []string
	Options           AuthorizationOptions
	ValidationOptions AuthorizationValidationOptions
	KeyVault          *jwt_keyvault.JwtKeyVaultService
	ContextAdapter    interfaces.UserDatabaseAdapter
}

var currentAuthorizationContext *AuthorizationContext

func NewFromUser(user *ContextUser) *AuthorizationContext {
	newContext := AuthorizationContext{
		User: user,
		ValidationOptions: AuthorizationValidationOptions{
			Audiences:     false,
			ExpiryDate:    true,
			Subject:       true,
			Issuer:        true,
			VerifiedEmail: false,
			NotBefore:     false,
			Tenant:        false,
		},
		Audiences: make([]string, 0),
	}

	newContext.KeyVault = jwt_keyvault.Get()
	newContext.WithDefaultOptions()

	currentAuthorizationContext = &newContext
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
	scope := config.GetString("JWT_SCOPE")
	authorizationType := config.GetString("JWT_AUTH_TYPE")
	keySize := config.GetString("JWT_KEY_SIZE")
	keyId := config.GetString("JWT_KEY_ID")
	if keyId != "" {
		keyId = "_" + keyId
	}
	privateKey := config.GetString("JWT" + keyId + "_PRIVATE_KEY")

	if issuer == "" {
		apiPort := config.GetString("HTTP_PORT")
		apiPrefix := config.GetString("API_PREFIX")
		issuer = "http://localhost"
		if apiPort != "" {
			issuer += ":" + apiPort
		}
		if apiPrefix != "" {
			if strings.HasPrefix(apiPrefix, "/") {
				issuer += apiPrefix
			} else {
				issuer += "/" + apiPrefix
			}
		}
		issuer += "/auth/global"
	}
	a.Issuer = issuer

	if tokenDuration == 0 {
		tokenDuration = 60
	}

	if scope == "" {
		scope = "authorization"
	}
	a.Scope = scope

	if authorizationType == "" {
		authorizationType = "hmac"
	}

	a.Options = AuthorizationOptions{
		TokenDuration: tokenDuration,
	}

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
	for _, inAudience := range a.Audiences {
		if strings.EqualFold(inAudience, audience) {
			found = true
			break
		}
	}
	if !found {
		a.Audiences = append(a.Audiences, audience)
	}

	return a
}

func (a *AuthorizationContext) WithIssuer(issuer string) *AuthorizationContext {
	a.Issuer = issuer

	return a
}

func (a *AuthorizationContext) WithDuration(tokenDuration int) *AuthorizationContext {
	a.Options.TokenDuration = tokenDuration

	return a
}

func (a *AuthorizationContext) WithScope(scope string) *AuthorizationContext {
	a.Scope = scope

	return a
}

func (a *AuthorizationContext) SetRequestIssuer(r *http.Request, tenantId string) string {
	baseUrl := service_provider.Get().GetBaseUrl(r)
	a.Issuer = baseUrl + "/auth/" + tenantId
	return a.Issuer
}

func GetCurrent() *AuthorizationContext {
	if currentAuthorizationContext != nil {
		return currentAuthorizationContext
	}

	return nil
}
