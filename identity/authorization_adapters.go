package identity

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/security"
	"github.com/pascaldekloe/jwt"
)

func AuthorizationAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.Get()
			authorized := true
			logger.Info("Authentication Layer Started")
			// Getting the token for validation
			jwt_token, valid := security.GetAuthorizationToken(r.Header)
			if !valid {
				authorized = false
			}

			// Validating token against the keys
			token, err := jwt.HMACCheck([]byte(jwt_token), []byte(security.PrivateKey))
			if err != nil {
				authorized = false
			}

			if authorized {
				user := executionctx.UserCtx{
					Name:      token.Subject,
					Audiences: token.Audiences,
					Issuer:    token.Issuer,
				}
				executionctx.NewUserContext(&user)
				logger.Info("User " + user.Name + " was authenticated successfully.")
			} else {
				response := LoginErrorResponse{
					Code:    "401",
					Error:   "Token is invalid",
					Message: "The token  is invalid",
				}

				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func EndAuthorizationAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.Get()
			ctx := executionctx.GetContext()
			if ctx.User != nil {
				logger.Info("Clearing user context from login")
				ctx.User = nil
			}
			next.ServeHTTP(w, r)
		})
	}
}

func EndAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Get()
		ctx := executionctx.GetContext()
		if ctx.User != nil {
			logger.Info("Clearing user context from login")
			ctx.User = nil
		}
		next.ServeHTTP(w, r)
	})
}
