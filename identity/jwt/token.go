// Package jwt provides the needed functions to generate tokens for users
package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/identity"
	"github.com/cjlapao/common-go/identity/jwt_keyvault"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security/encryption"
	"github.com/google/uuid"
	"github.com/pascaldekloe/jwt"
)

type RawCertificateHeader struct {
	X5T string `json:"x5t"`
}

// GenerateDefaultUserToken generates a jwt user token with the default audiences in the context
// It returns a user token object and an error if it exists
func GenerateDefaultUserToken(user models.User) (*models.UserToken, error) {
	ctx := execution_context.Get()

	return GenerateUserTokenForAudiences("", user, ctx.Authorization.Options.Audiences...)
}

// GenerateUserToken
func GenerateUserToken(keyId string, user models.User) (*models.UserToken, error) {
	ctx := execution_context.Get()

	return GenerateUserTokenForAudiences(keyId, user, ctx.Authorization.Options.Audiences...)
}

func GenerateUserTokenForAudiences(keyId string, user models.User, audiences ...string) (*models.UserToken, error) {
	ctx := execution_context.Get()

	return GenerateUserTokenForKeyAndAudiences(keyId, user, ctx.Authorization.Options.Audiences...)
}

func GenerateUserTokenForKeyAndAudiences(keyId string, user models.User, audiences ...string) (*models.UserToken, error) {
	var userToken models.UserToken
	var userTokenClaims jwt.Claims
	ctx := execution_context.Get()
	now := time.Now().Round(time.Second)
	nowSkew := now.Add((time.Minute * 2))
	nowNegativeSkew := now.Add((time.Minute * 2) * -1)
	validUntil := nowSkew.Add(time.Hour * 1)

	userTokenClaims.Subject = user.Email
	userTokenClaims.Issuer = ctx.Authorization.Options.Issuer
	userTokenClaims.Issued = jwt.NewNumericTime(nowSkew)
	userTokenClaims.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	userTokenClaims.Expires = jwt.NewNumericTime(validUntil)
	userTokenClaims.ID = uuid.NewString()

	// Adding Custom Claims to the token
	userClaims := make(map[string]interface{})
	userClaims["scope"] = ctx.Authorization.Options.Scope
	userClaims["uid"] = strings.ToLower(user.ID)
	userClaims["name"] = user.DisplayName
	userClaims["given_name"] = user.FirstName
	userClaims["family_name"] = user.LastName

	// Adding the email verification to the token if the validation is on
	if ctx.Authorization.ValidationOptions.VerifiedEmail {
		userClaims["email_verified"] = false
	}

	// Adding the correlation nonce to the token if it exists
	if ctx.Authorization.CorrelationId != "" {
		userClaims["nonce"] = ctx.CorrelationId
	}

	// Adding the tenantId if it exists
	if ctx.Authorization.TenantId != "" {
		userClaims["tid"] = ctx.Authorization.TenantId
	}

	userTokenClaims.KeyID = ctx.Authorization.Options.KeyId

	// Reading all of the roles
	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	userClaims["roles"] = roles

	// Reading all the audiences
	if len(audiences) > 0 {
		userTokenClaims.Audiences = audiences
	}

	userTokenClaims.Set = userClaims

	var token string
	var err error

	token, err = signToken(keyId, userTokenClaims)
	if err != nil {
		logger.Error("There was an error generating a jwt token for user %v with key id %v", user.Username, keyId)
		return nil, err
	}

	userToken = models.UserToken{
		Token:     token,
		ExpiresAt: validUntil,
		NotBefore: nowNegativeSkew,
		Audiences: audiences,
		Issuer:    userTokenClaims.Issuer,
		UsedKeyID: keyId,
	}

	refreshToken, err := GenerateRefreshToken(keyId, user)
	if err != nil {
		userToken.RefreshToken = refreshToken
	}

	return &userToken, nil
}

// GenerateRefreshToken generates a refresh token for the user with a
func GenerateRefreshToken(keyId string, user models.User) (string, error) {
	var refreshTokenClaims jwt.Claims
	ctx := execution_context.Get()
	now := time.Now().Round(time.Second)
	nowSkew := now.Add((time.Hour * 2))
	nowNegativeSkew := now.Add((time.Minute * 2) * -1)
	validUntil := nowSkew.Add((time.Hour * 24) * 365)

	refreshTokenClaims.Subject = user.Email
	refreshTokenClaims.Issuer = ctx.Authorization.Options.Issuer
	refreshTokenClaims.Issued = jwt.NewNumericTime(nowSkew)
	refreshTokenClaims.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	refreshTokenClaims.Expires = jwt.NewNumericTime(validUntil)
	refreshTokenClaims.ID = uuid.NewString()

	// Custom Claims
	userClaims := make(map[string]interface{})
	userClaims["scope"] = identity.RefreshTokenScope
	userClaims["name"] = user.DisplayName
	userClaims["given_name"] = user.FirstName
	userClaims["family_name"] = user.LastName
	userClaims["uid"] = strings.ToLower(user.ID)
	if ctx.Authorization.TenantId != "" {
		userClaims["tid"] = ctx.Authorization.TenantId
	}
	refreshTokenClaims.KeyID = ctx.Authorization.Options.KeyId

	refreshToken, err := signToken(keyId, refreshTokenClaims)
	if err != nil {
		logger.Error("There was an error signing the refresh token for user %v with key id %v", user.Username, keyId)
		return "", err
	}

	return refreshToken, nil
}

func GenerateVerifyEmailToken(keyId string, user models.User) string {
	var refreshTokenClaims jwt.Claims
	ctx := execution_context.Get()
	now := time.Now().Round(time.Second)
	nowSkew := now.Add((time.Hour * 2))
	nowNegativeSkew := now.Add((time.Minute * 2) * -1)
	validUntil := nowSkew.Add(time.Hour * 2)

	refreshTokenClaims.Subject = user.Email
	refreshTokenClaims.Issuer = ctx.Authorization.Options.Issuer
	refreshTokenClaims.Issued = jwt.NewNumericTime(nowSkew)
	refreshTokenClaims.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	refreshTokenClaims.Expires = jwt.NewNumericTime(validUntil)
	refreshTokenClaims.ID = uuid.NewString()

	// Custom Claims
	userClaims := make(map[string]interface{})
	userClaims["scope"] = identity.EmailVerificationScope
	userClaims["name"] = user.DisplayName
	userClaims["given_name"] = user.FirstName
	userClaims["family_name"] = user.LastName
	userClaims["uid"] = strings.ToLower(user.ID)
	userClaims["tid"] = ctx.Authorization.TenantId
	refreshTokenClaims.KeyID = ctx.Authorization.Options.KeyId

	refreshToken, err := signToken(keyId, refreshTokenClaims)
	if err != nil {
		logger.Error("There was an error signing the email verification token for user %v with key id %v", user.Username, keyId)
		return ""
	}

	return refreshToken
}

func signToken(keyId string, claims jwt.Claims) (string, error) {
	ctx := execution_context.Get()
	var rawToken []byte
	var err error
	var signKey *jwt_keyvault.JwtKeyVaultItem
	if keyId == "" {
		signKey = ctx.Authorization.KeyVault.GetDefaultKey()
	} else {
		signKey = ctx.Authorization.KeyVault.GetKey(keyId)
	}
	if signKey == nil {
		err = errors.New("signing key was not found")
		logger.Error("There was an error signing the token with key %v, it was not found int the key vault", keyId)
		return "", err
	}

	// Adding extra headers for some signing cases
	extraHeaders, _ := json.Marshal(RawCertificateHeader{
		X5T: signKey.ID,
	})

	// Signing the token using the key encryption type
	switch kt := signKey.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		switch signKey.Size {
		case encryption.Bit256:
			rawToken, err = claims.ECDSASign("ES256", kt, extraHeaders)
		case encryption.Bit384:
			rawToken, err = claims.ECDSASign("ES384", kt, extraHeaders)
		case encryption.Bit512:
			rawToken, err = claims.ECDSASign("ES512", kt, extraHeaders)
		}
	case string:
		switch signKey.Size {
		case encryption.Bit256:
			rawToken, err = claims.HMACSign("HS256", []byte(kt))
		case encryption.Bit384:
			rawToken, err = claims.HMACSign("HS384", []byte(kt))
		case encryption.Bit512:
			rawToken, err = claims.HMACSign("HS512", []byte(kt))
		}
	case *rsa.PrivateKey:
		switch signKey.Size {
		case encryption.Bit256:
			rawToken, err = claims.RSASign("RS256", kt, extraHeaders)
		case encryption.Bit384:
			rawToken, err = claims.RSASign("RS384", kt, extraHeaders)
		case encryption.Bit512:
			rawToken, err = claims.RSASign("RS512", kt, extraHeaders)
		}
	}

	if err != nil {
		return "", err
	}

	return string(rawToken), nil
}
