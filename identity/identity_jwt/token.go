package identity_jwt

import (
	"time"

	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security"
	"github.com/pascaldekloe/jwt"
)

// GenerateUserToken generates a jwt user token
func GenerateUserToken(user models.User) (string, string) {
	ctx := execution_context.Get()
	var userToken jwt.Claims
	now := time.Now().Round(time.Second)
	nowSkew := now.Add((time.Minute * 2))
	nowNegativeSkew := now.Add((time.Minute * 2) * -1)
	validUntil := nowSkew.Add(time.Hour * 1)

	userToken.Subject = user.Email
	userToken.Issuer = ctx.Authorization.Options.Issuer
	userToken.Issued = jwt.NewNumericTime(nowSkew)
	userToken.NotBefore = jwt.NewNumericTime(nowNegativeSkew)
	userToken.Expires = jwt.NewNumericTime(validUntil)
	userClaims := make(map[string]interface{})
	userClaims["email_verified"] = false
	userClaims["scope"] = ctx.Authorization.Options.Scope
	userClaims["name"] = user.DisplayName
	userClaims["given_name"] = user.FirstName
	userClaims["nonce"] = ctx.CorrelationId
	userClaims["family_name"] = user.LastName
	roles := make([]string, 0)
	for _, role := range user.Roles {
		roles = append(roles, role.Name)
	}
	userClaims["roles"] = roles

	userToken.Set = userClaims

	var token []byte
	var err error

	signType := ctx.Authorization.Options.SignatureType
	signSize := ctx.Authorization.Options.SignatureSize
	userToken.KeyID = ctx.Authorization.Options.KeyId
	switch signType {
	case security.ECDSA:
		switch signSize {
		case security.Bit256:
			privateKey, _ := security.ECDSAHelper{}.Decode(ctx.Authorization.Options.PrivateKey, ctx.Authorization.Options.PublicKey)
			token, err = userToken.ECDSASign("ES256", privateKey)
		case security.Bit384:
			privateKey, _ := security.ECDSAHelper{}.Decode(ctx.Authorization.Options.PrivateKey, ctx.Authorization.Options.PublicKey)
			token, err = userToken.ECDSASign("ES384", privateKey)
		case security.Bit512:
			privateKey, _ := security.ECDSAHelper{}.Decode(ctx.Authorization.Options.PrivateKey, ctx.Authorization.Options.PublicKey)
			token, err = userToken.ECDSASign("ES512", privateKey)
		}
	case security.HMAC:
		switch signSize {
		case security.Bit256:
			token, err = userToken.HMACSign("HS256", []byte(ctx.Authorization.Options.PrivateKey))
		case security.Bit384:
			token, err = userToken.HMACSign("HS384", []byte(ctx.Authorization.Options.PrivateKey))
		case security.Bit512:
			token, err = userToken.HMACSign("HS512", []byte(ctx.Authorization.Options.PrivateKey))
		}
	}

	helper.CheckError(err)

	return string(token), userToken.Expires.String()
}
