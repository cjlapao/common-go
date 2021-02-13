package security

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/cjlapao/common-go/helper"

	loggerModule "github.com/cjlapao/common-go/log"

	"github.com/pascaldekloe/jwt"
)

var logger = loggerModule.Get()

// Security Constants
const (
	Issuer     = "Ittech24.co.uk"
	LoginScope = "authorization"
	PrivateKey = "somerandomshit"
)

// SHA256Encode Hash string with SHA256
func SHA256Encode(value string) string {
	hasher := sha256.New()
	bytes := []byte(value)
	hasher.Write(bytes)

	return hex.EncodeToString(hasher.Sum(nil))
}

// GenerateUserToken generates a jwt user token
func GenerateUserToken(email string) (string, string) {

	var claims jwt.Claims

	claims.Subject = email
	claims.Issuer = Issuer
	claims.Issued = jwt.NewNumericTime(time.Now().Round(time.Second))
	claims.Expires = jwt.NewNumericTime(time.Now().Add(time.Hour * 1))
	claims.Set = map[string]interface{}{"email_verified": false, "scope": "authentication"}

	token, err := claims.HMACSign("HS256", []byte(PrivateKey))

	helper.CheckError(err)

	return string(token), claims.Expires.String()
}

func ValidateToken(token string) bool {
	claims, err := jwt.HMACCheck([]byte(token), []byte(PrivateKey))

	if err != nil {
		logger.Error("Token is not valid ")
		return false
	}
	email := claims.Subject
	if !claims.Valid(time.Now()) {
		logger.Error("Token is not valid for user " + email)
		return false
	}

	if claims.Issuer != Issuer {
		logger.Error("Token is not valid for user " + email)
		return false
	}

	return true
}

func GetAuthorizationToken(request http.Header) (string, bool) {
	authHeader := strings.Split(request.Get("Authorization"), "Bearer ")
	if len(authHeader) != 2 {
		return "", false
	}

	logger.Debug("Token: " + authHeader[1])

	return authHeader[1], true
}

func AuthenticateMiddleware(target http.HandlerFunc) http.Handler {
	next := http.Handler(target)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, valid := GetAuthorizationToken(r.Header)
		if !valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		if !ValidateToken(token) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
