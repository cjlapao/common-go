package middleware

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity/authorization_context"
	"github.com/cjlapao/common-go/identity/jwt"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/log"
	"github.com/gorilla/mux"
)

func AuthorizationMiddlewareAdapter(roles []string, claims []string) controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := execution_context.Get()
			vars := mux.Vars(r)
			tenantId := vars["tenantId"]

			// if no tenant is set we will assume it is the global tenant
			if tenantId == "" {
				tenantId = "global"
			}

			ctx.Authorization.SetRequestIssuer(r, tenantId)
			authorized := true
			logger.Info("Authorization Layer Started")

			for _, claim := range claims {
				println(claim)
			}

			for _, role := range roles {
				println(role)
			}

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
				// Getting user roles and claims
				dbUser := ctx.Authorization.ContextAdapter.GetUserByEmail(userToken.User)
				validatedRoles := make(map[string]bool)
				// validatedClaims := make(map[string]bool)
				if !dbUser.IsValid() {
					authorized = false
				} else {
					if len(roles) > 0 {
						if dbUser.Roles == nil || len(dbUser.Roles) == 0 {
							authorized = false
						} else {
							for _, requiredRole := range roles {
								foundRole := false
								for _, role := range dbUser.Roles {
									if strings.EqualFold(requiredRole, role.ID) {
										foundRole = true
										break
									}
								}

								validatedRoles[requiredRole] = foundRole
							}
						}
					}
				}

				for _, found := range validatedRoles {
					if !found {
						authorized = false
						break
					}
				}
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
