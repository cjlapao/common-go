package middleware

import (
	"encoding/json"
	"errors"
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

var ctx = execution_context.Get()

func AuthorizationMiddlewareAdapter(roles []string, claims []string) controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			tenantId := vars["tenantId"]
			var userToken *models.UserToken
			var dbUser *models.User
			var validateError error
			var err error

			// if no tenant is set we will assume it is the global tenant
			if tenantId == "" {
				tenantId = "global"
			}

			// Setting the tenant in the context
			ctx.Authorization.SetRequestIssuer(r, tenantId)

			//Starting authorization layer of the token
			authorized := true
			logger.Info("Authorization Layer Started")

			// Getting the token for validation
			jwt_token, valid := http_helper.GetAuthorizationToken(r.Header)
			if !valid {
				authorized = false
				validateError = errors.New("bearer token not found in request")
				logger.Error(validateError.Error())
			}

			// Validating userToken against the keys
			if authorized {
				userToken, err = jwt.ValidateUserToken(jwt_token, ctx.Authorization.Options.Scope)
				if validateError != nil {
					authorized = false
					validateError = errors.New("bearer token is not valid, " + err.Error())
					logger.Error(validateError.Error())
				}
			}

			// Making sure the user token does contain the necessary fields
			if userToken == nil || userToken.User == "" {
				authorized = false
				validateError = errors.New("bearer token has invalid user")
				logger.Error(validateError.Error())
			}

			// To gain speed we will only get the db user if there is any role or claim to validate
			// otherwise we don't need anything else to validate it
			if len(roles) > 0 || len(claims) > 0 {
				// Getting the user from the database to validate roles and claims
				if authorized {
					dbUser = ctx.UserDatabaseAdapter.GetUserByEmail(userToken.User)
					if dbUser == nil || dbUser.ID == "" {
						authorized = false
						validateError = errors.New("bearer token user was not found in database, potentially revoked, " + userToken.User)
						logger.Error(validateError.Error())
					}
				}

				// Validating user roles
				if authorized {
					err = validateUserRoles(dbUser, roles)
					if err != nil {
						authorized = false
						validateError = errors.New("bearer token user does not contain one or more roles required by the context, " + err.Error())
						logger.Error(validateError.Error())
					}
				}

				// Validating user claims
				if authorized {
					err = validateUserClaims(dbUser, claims)
					if err != nil {
						authorized = false
						validateError = errors.New("bearer token user does not contain one or more claims required by the context, " + err.Error())
						logger.Error(validateError.Error())
					}
				}
			}

			if authorized && userToken != nil && userToken.ID != "" {
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
					ErrorDescription: validateError.Error(),
				}

				w.WriteHeader(http.StatusUnauthorized)
				json.NewEncoder(w).Encode(response)
				if userToken != nil && userToken.User != "" {
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

func validateUserRoles(user *models.User, requiredRoles []string) error {
	var validateError error
	if user == nil || user.ID == "" {
		validateError = errors.New("user is invalid when processing roles")
		logger.Error(validateError.Error())
		return validateError
	}

	// Getting user roles and claims
	validatedRoles := make(map[string]bool)

	if !user.IsValid() {
		validateError = errors.New("user is invalid when processing roles")
		logger.Error(validateError.Error())
		return validateError
	} else {
		if len(requiredRoles) > 0 {
			if user.Roles == nil || len(user.Roles) == 0 {
				validateError = errors.New("user does not contain any roles")
				logger.Error(validateError.Error())
				return validateError
			} else {
				for _, requiredRole := range requiredRoles {
					foundRole := false
					for _, role := range user.Roles {
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

	for roleName, found := range validatedRoles {
		if !found {
			validateError = errors.New("user does not contain required role " + roleName)
			logger.Error(validateError.Error())
			return validateError
		}
	}

	return nil
}

func validateUserClaims(user *models.User, requiredClaims []string) error {
	var validateError error
	if user == nil || user.ID == "" {
		validateError = errors.New("user is invalid when processing claims")
		logger.Error(validateError.Error())
		return validateError
	}

	// Getting user claims and claims
	validatedClaims := make(map[string]bool)

	if !user.IsValid() {
		validateError = errors.New("user is invalid when processing claims")
		logger.Error(validateError.Error())
		return validateError
	} else {
		if len(requiredClaims) > 0 {
			if user.Claims == nil || len(user.Claims) == 0 {
				validateError = errors.New("user does not contain any claims")
				logger.Error(validateError.Error())
				return validateError
			} else {
				for _, requiredClaim := range requiredClaims {
					foundClaim := false
					for _, claim := range user.Claims {
						if strings.EqualFold(requiredClaim, claim.ID) {
							foundClaim = true
							break
						}
					}

					validatedClaims[requiredClaim] = foundClaim
				}
			}
		}
	}

	for claimName, found := range validatedClaims {
		if !found {
			validateError = errors.New("user does not contain required claim " + claimName)
			logger.Error(validateError.Error())
			return validateError
		}
	}

	return nil
}
