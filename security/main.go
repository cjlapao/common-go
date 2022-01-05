package security

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/log"
	"github.com/pascaldekloe/jwt"
)

var logger = log.Get()

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

func AuthenticateMiddleware(target http.HandlerFunc) http.Handler {
	next := http.Handler(target)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, valid := http_helper.GetAuthorizationToken(r.Header)
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
