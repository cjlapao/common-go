package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity/authorization_context"
	"github.com/cjlapao/common-go/identity/jwt"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/log"
)

func AuthorizationMiddlewareAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := execution_context.Get()
			authorized := true
			logger.Info("Authorization Layer Started")

			// Getting the token for validation
			jwt_token, valid := http_helper.GetAuthorizationToken(r.Header)
			if !valid {
				authorized = false
			}

			// Validating userToken against the keys
			userToken, err := jwt.ValidateUserToken(jwt_token, ctx.Authorization.Options.Scope)
			if err != nil {
				authorized = false
			}

			if authorized {
				user := authorization_context.NewContextUser()
				user.ID = userToken.ID
				user.Email = userToken.User
				user.Audiences = userToken.Audiences
				user.Issuer = userToken.Issuer

				ctx.Authorization = authorization_context.NewFromUser(user)
				ctx.Authorization.TenantId = userToken.TenantId
				logger.Info("User " + user.Email + " was authorized successfully.")
			} else {
				response := models.OAuthErrorResponse{
					Error:            models.OAuthUnauthorizedClient,
					ErrorDescription: "The token is invalid",
				}

				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				if userToken != nil {
					logger.Error("User " + userToken.User + " failed to authorize.")
				} else {
					logger.Error("Request failed to authorize. No bearer token found.")

				}

				logger.Info("Authorization Layer Finished")
				return
			}

			logger.Info("Authorization Layer Finished")
			next.ServeHTTP(w, r)
		})
	}
}

func EndAuthorizationMiddlewareAdapter() controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger := log.Get()
			ctx := execution_context.Get()
			authCtx := authorization_context.GetCurrent()
			if authCtx.User != nil {
				logger.Info("Clearing user context from login")
				authCtx.User = nil
				ctx.CorrelationId = ""
			}
			next.ServeHTTP(w, r)
		})
	}
}
