// Package jwt provides the needed functions to generate tokens for users
package jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cjlapao/common-go/cryptorand"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/identity"
	"github.com/cjlapao/common-go/identity/jwt_keyvault"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security/encryption"
	"github.com/pascaldekloe/jwt"
)

func GetTokenClaim(token string, claim string) string {
	if token == "" || claim == "" {
		return ""
	}

	jwtToken, err := jwt.ParseWithoutCheck([]byte(token))

	if err != nil {
		return ""
	}

	// Transforming token into a user token
	rawJsonToken, _ := jwtToken.Raw.MarshalJSON()
	var tokenMap map[string]interface{}
	err = json.Unmarshal(rawJsonToken, &tokenMap)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%v", tokenMap[claim])
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
	if ctx.Authorization.ValidationOptions.NotBefore {
		userTokenClaims.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	}
	userTokenClaims.Expires = jwt.NewNumericTime(validUntil)
	userTokenClaims.ID = cryptorand.GenerateAlphaNumericRandomString(60)

	// Adding Custom Claims to the token
	userClaims := make(map[string]interface{})
	userClaims["scope"] = ctx.Authorization.Options.Scope
	userClaims["uid"] = strings.ToLower(user.ID)
	userClaims["name"] = user.DisplayName
	userClaims["given_name"] = user.FirstName
	userClaims["family_name"] = user.LastName

	// Adding the email verification to the token if the validation is on
	if ctx.Authorization.ValidationOptions.VerifiedEmail {
		userClaims["email_verified"] = user.EmailVerified
	}

	// Adding the correlation nonce to the token if it exists
	if ctx.CorrelationId != "" {
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
	if err == nil {
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
	if ctx.Authorization.ValidationOptions.NotBefore {
		refreshTokenClaims.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	}
	refreshTokenClaims.Expires = jwt.NewNumericTime(validUntil)
	refreshTokenClaims.ID = cryptorand.GenerateAlphaNumericRandomString(60)

	// Custom Claims
	customClaims := make(map[string]interface{})
	customClaims["scope"] = identity.RefreshTokenScope
	customClaims["name"] = user.DisplayName
	customClaims["given_name"] = user.FirstName
	customClaims["family_name"] = user.LastName
	customClaims["uid"] = strings.ToLower(user.ID)
	if ctx.Authorization.TenantId != "" {
		customClaims["tid"] = ctx.Authorization.TenantId
	}
	refreshTokenClaims.KeyID = ctx.Authorization.Options.KeyId
	refreshTokenClaims.Set = customClaims

	refreshToken, err := signToken(keyId, refreshTokenClaims)
	if err != nil {
		logger.Error("There was an error signing the refresh token for user %v with key id %v", user.Username, keyId)
		return "", err
	}

	return refreshToken, nil
}

func GenerateVerifyEmailToken(keyId string, user models.User) string {
	var emailVerificationTokenClaims jwt.Claims
	ctx := execution_context.Get()
	now := time.Now().Round(time.Second)
	nowSkew := now.Add((time.Hour * 2))
	nowNegativeSkew := now.Add((time.Minute * 2) * -1)
	validUntil := nowSkew.Add(time.Hour * 2)

	emailVerificationTokenClaims.Subject = user.Email
	emailVerificationTokenClaims.Issuer = ctx.Authorization.Options.Issuer
	emailVerificationTokenClaims.Issued = jwt.NewNumericTime(nowSkew)
	if ctx.Authorization.ValidationOptions.NotBefore {
		emailVerificationTokenClaims.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	}
	emailVerificationTokenClaims.Expires = jwt.NewNumericTime(validUntil)
	emailVerificationTokenClaims.ID = cryptorand.GenerateAlphaNumericRandomString(60)

	// Custom Claims
	customClaims := make(map[string]interface{})
	customClaims["scope"] = identity.EmailVerificationScope
	customClaims["name"] = user.DisplayName
	customClaims["given_name"] = user.FirstName
	customClaims["family_name"] = user.LastName
	customClaims["uid"] = strings.ToLower(user.ID)
	if ctx.Authorization.TenantId != "" {
		customClaims["tid"] = ctx.Authorization.TenantId
	}
	emailVerificationTokenClaims.KeyID = ctx.Authorization.Options.KeyId
	emailVerificationTokenClaims.Set = customClaims
	refreshToken, err := signToken(keyId, emailVerificationTokenClaims)
	if err != nil {
		logger.Error("There was an error signing the email verification token for user %v with key id %v", user.Username, keyId)
		return ""
	}

	return refreshToken
}

func ValidateUserToken(token string, scope string, audiences ...string) (*models.UserToken, error) {
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	ctx := execution_context.Get()
	var tokenBytes []byte
	var verifiedToken *jwt.Claims
	tokenBytes = []byte(token)
	var err error
	var signKey *jwt_keyvault.JwtKeyVaultItem
	rawToken, err := jwt.ParseWithoutCheck(tokenBytes)
	if err != nil {
		return nil, err
	}

	// Verifying signature using the key that was sign with
	signKey = ctx.Authorization.KeyVault.GetKey(rawToken.KeyID)
	switch kt := signKey.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		key := kt.PublicKey
		verifiedToken, err = jwt.ECDSACheck(tokenBytes, &key)
		if err != nil {
			return nil, err
		}
	case string:
		verifiedToken, err = jwt.HMACCheck(tokenBytes, []byte(kt))
		if err != nil {
			return nil, err
		}
	case *rsa.PrivateKey:
		key := kt.PublicKey
		verifiedToken, err = jwt.RSACheck(tokenBytes, &key)
		if err != nil {
			return nil, err
		}
	}

	// Transforming token into a user token
	rawJsonToken, _ := verifiedToken.Raw.MarshalJSON()
	var userToken models.UserToken
	err = json.Unmarshal(rawJsonToken, &userToken)
	if err != nil {
		return nil, errors.New("token is not formated correctly")
	}

	// Validating the scope of the token
	if !strings.EqualFold(scope, userToken.Scope) {
		return &userToken, errors.New("token scope is not valid")
	}

	// Validating the token not before property
	if ctx.Authorization.ValidationOptions.NotBefore {
		if userToken.NotBefore.After(time.Now()) {
			return &userToken, errors.New("token is not yet valid")
		}
	}

	// Validating expiry token
	if userToken.ExpiresAt.Before(time.Now()) {
		return &userToken, errors.New("token is expired")
	}

	// If we require the Issuer to be validated we will be validating it
	if ctx.Authorization.ValidationOptions.Issuer {
		if !strings.EqualFold(userToken.Issuer, ctx.Authorization.Options.Issuer) {
			return &userToken, errors.New("token is not valid for subject " + userToken.DisplayName)
		}
	}

	// Validating if the email has been verified
	if ctx.Authorization.ValidationOptions.VerifiedEmail {
		if !userToken.EmailVerified {
			return &userToken, errors.New("email is not verified for subject " + userToken.DisplayName)
		}
	}

	// Validating if the token contains the necessary audiences
	if ctx.Authorization.ValidationOptions.Audiences && len(audiences) > 0 {
		if len(audiences) == 0 || len(userToken.Audiences) == 0 {
			return &userToken, errors.New("no audiences to validate subject " + userToken.DisplayName)
		}
		isValid := true
		for _, audience := range audiences {
			wasFound := false
			for _, userAudience := range userToken.Audiences {
				if strings.EqualFold(userAudience, audience) {
					wasFound = true
				}
			}
			if !wasFound {
				isValid = false
				break
			}
		}

		if !isValid {
			return &userToken, errors.New("one or more required audience was not found for subject " + userToken.DisplayName)
		}
	}

	// Validating if the token tenant id is the same as the context
	if ctx.Authorization.ValidationOptions.Tenant {
		if ctx.Authorization.TenantId == "" || userToken.TenantId == "" {
			return &userToken, errors.New("no tenant was not found for subject " + userToken.DisplayName)
		}
		if !strings.EqualFold(ctx.Authorization.TenantId, userToken.TenantId) {
			return &userToken, errors.New("token is not valid for tenant " + userToken.TenantId + " for subject " + userToken.DisplayName)
		}
	}

	return &userToken, nil
}

func ValidateRefreshToken(token string, user string) (*models.UserToken, error) {
	if token == "" {
		return nil, errors.New("token cannot be empty")
	}

	ctx := execution_context.Get()
	var tokenBytes []byte
	var verifiedToken *jwt.Claims
	tokenBytes = []byte(token)
	var err error
	var signKey *jwt_keyvault.JwtKeyVaultItem
	rawToken, err := jwt.ParseWithoutCheck(tokenBytes)
	if err != nil {
		return nil, err
	}

	// Verifying signature using the key that was sign with
	signKey = ctx.Authorization.KeyVault.GetKey(rawToken.KeyID)
	switch kt := signKey.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		key := kt.PublicKey
		verifiedToken, err = jwt.ECDSACheck(tokenBytes, &key)
		if err != nil {
			return nil, err
		}
	case string:
		verifiedToken, err = jwt.HMACCheck(tokenBytes, []byte(kt))
		if err != nil {
			return nil, err
		}
	case *rsa.PrivateKey:
		key := kt.PublicKey
		verifiedToken, err = jwt.RSACheck(tokenBytes, &key)
		if err != nil {
			return nil, err
		}
	}

	// Transforming token into a user token
	rawJsonToken, _ := verifiedToken.Raw.MarshalJSON()
	var userToken models.UserToken
	err = json.Unmarshal(rawJsonToken, &userToken)
	if err != nil {
		return nil, errors.New("token is not formated correctly")
	}

	// Validating the scope of the token
	if !strings.EqualFold(user, userToken.User) {
		return &userToken, errors.New("token user is not valid")
	}

	// Validating the scope of the token
	if !strings.EqualFold(identity.RefreshTokenScope, userToken.Scope) {
		return &userToken, errors.New("token scope is not valid")
	}

	// Validating expiry token
	if userToken.ExpiresAt.Before(time.Now()) {
		return &userToken, errors.New("token is expired")
	}

	// If we require the Issuer to be validated we will be validating it
	if ctx.Authorization.ValidationOptions.Issuer {
		if !strings.EqualFold(userToken.Issuer, ctx.Authorization.Options.Issuer) {
			return &userToken, errors.New("token is not valid for subject " + userToken.DisplayName)
		}
	}

	// Validating if the token tenant id is the same as the context
	if ctx.Authorization.ValidationOptions.Tenant {
		if ctx.Authorization.TenantId == "" || userToken.TenantId == "" {
			return &userToken, errors.New("no tenant was not found for subject " + userToken.DisplayName)
		}
		if !strings.EqualFold(ctx.Authorization.TenantId, userToken.TenantId) {
			return &userToken, errors.New("token is not valid for tenant " + userToken.TenantId + " for subject " + userToken.DisplayName)
		}
	}

	return &userToken, nil
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
	var extraHeaders []byte

	// Signing the token using the key encryption type
	switch kt := signKey.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		// Adding extra headers for some signing cases
		extraHeaders, _ = json.Marshal(RawCertificateHeader{
			KeyId: signKey.ID,
			X5T:   signKey.Thumbprint,
		})
		switch signKey.Size {
		case encryption.Bit256:
			rawToken, err = claims.ECDSASign("ES256", kt, extraHeaders)
		case encryption.Bit384:
			rawToken, err = claims.ECDSASign("ES384", kt, extraHeaders)
		case encryption.Bit512:
			rawToken, err = claims.ECDSASign("ES512", kt, extraHeaders)
		}
	case string:
		// Adding extra headers for some signing cases
		extraHeaders, _ = json.Marshal(RawCertificateHeader{
			KeyId: signKey.ID,
		})
		switch signKey.Size {
		case encryption.Bit256:
			rawToken, err = claims.HMACSign("HS256", []byte(kt), extraHeaders)
		case encryption.Bit384:
			rawToken, err = claims.HMACSign("HS384", []byte(kt), extraHeaders)
		case encryption.Bit512:
			rawToken, err = claims.HMACSign("HS512", []byte(kt), extraHeaders)
		}
	case *rsa.PrivateKey:
		// Adding extra headers for some signing cases
		extraHeaders, _ = json.Marshal(RawCertificateHeader{
			KeyId: signKey.ID,
			X5T:   signKey.Thumbprint,
		})
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
