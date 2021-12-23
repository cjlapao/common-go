package identity_jwt

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"encoding/json"
	"strings"
	"time"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security/encryption"
	"github.com/google/uuid"
	"github.com/pascaldekloe/jwt"
)

type RawCertificateHeader struct {
	X5T string `json:"x5t"`
}

// GenerateUserToken generates a jwt user token
func GenerateUserToken(user models.User) (string, string) {
	ctx := execution_context.Get()

	return GenerateUserTokenForAudiences(user, ctx.Authorization.User.Audiences...)
}

func GenerateUserTokenForAudiences(user models.User, audiences ...string) (string, string) {
	ctx := execution_context.Get()
	var userToken jwt.Claims
	now := time.Now().Round(time.Second)
	nowSkew := now.Add((time.Minute * 2))
	nowNegativeSkew := now.Add((time.Minute * 2) * -1)
	validUntil := nowSkew.Add(time.Hour * 1)

	userToken.Subject = user.Email
	userToken.Issuer = ctx.Authorization.Options.Issuer
	userToken.Audiences = ctx.Authorization.Options.Audiences
	userToken.Issued = jwt.NewNumericTime(nowSkew)
	userToken.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	userToken.Expires = jwt.NewNumericTime(validUntil)
	userToken.ID = uuid.NewString()

	// Custom Claims
	userClaims := make(map[string]interface{})
	userClaims["email_verified"] = false
	userClaims["scope"] = ctx.Authorization.Options.Scope
	userClaims["name"] = user.DisplayName
	userClaims["given_name"] = user.FirstName
	userClaims["nonce"] = ctx.CorrelationId
	userClaims["family_name"] = user.LastName
	userClaims["uid"] = strings.ToLower(user.ID)
	userClaims["tid"] = ctx.Authorization.TenantId
	userToken.KeyID = ctx.Authorization.Options.KeyId

	// Reading all of the roles
	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	userClaims["roles"] = roles

	userToken.Set = userClaims

	var token []byte
	var err error

	defaultKey := ctx.Authorization.KeyVault.GetDefaultKey()
	extraHeaders, _ := json.Marshal(RawCertificateHeader{
		X5T: defaultKey.ID,
	})

	switch kt := defaultKey.PrivateKey.(type) {
	case *ecdsa.PrivateKey:
		switch defaultKey.Size {
		case encryption.Bit256:
			token, err = userToken.ECDSASign("ES256", kt, extraHeaders)
		case encryption.Bit384:
			token, err = userToken.ECDSASign("ES384", kt, extraHeaders)
		case encryption.Bit512:
			token, err = userToken.ECDSASign("ES512", kt, extraHeaders)
		}
	case string:
		switch defaultKey.Size {
		case encryption.Bit256:
			token, err = userToken.HMACSign("HS256", []byte(kt))
		case encryption.Bit384:
			token, err = userToken.HMACSign("HS384", []byte(kt))
		case encryption.Bit512:
			token, err = userToken.HMACSign("HS512", []byte(kt))
		}
	case *rsa.PrivateKey:
		switch defaultKey.Size {
		case encryption.Bit256:
			token, err = userToken.RSASign("RS256", kt, extraHeaders)
		case encryption.Bit384:
			token, err = userToken.RSASign("RS384", kt, extraHeaders)
		case encryption.Bit512:
			token, err = userToken.RSASign("RS512", kt, extraHeaders)
		}
	}

	helper.CheckError(err)

	return string(token), userToken.Expires.String()
}
