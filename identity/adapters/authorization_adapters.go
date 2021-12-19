package adapters

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/identity/authorization_context"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/security"
	"github.com/pascaldekloe/jwt"
)

func AuthorizationAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := execution_context.Get()
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
				user := authorization_context.NewContextUser()
				user.ID = token.ID
				user.Email = token.Subject
				user.Audiences = token.Audiences
				user.Issuer = token.Issuer

				ctx.Authorization = authorization_context.NewFromUser(user)
				logger.Info("User " + user.DisplayName + " was authenticated successfully.")
			} else {
				response := models.LoginErrorResponse{
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
			ctx := authorization_context.GetCurrent()
			if ctx.User != nil {
				logger.Info("Clearing user context from login")
				ctx.User = nil
				ctx.CorrelationId = ""
			}
			next.ServeHTTP(w, r)
		})
	}
}

func EndAuthorizationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := log.Get()
		ctx := authorization_context.GetCurrent()
		if ctx.User != nil {
			logger.Info("Clearing user context from login")
			ctx.User = nil
			ctx.CorrelationId = ""
		}
		next.ServeHTTP(w, r)
	})
}
