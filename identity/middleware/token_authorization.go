package middleware

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity"
	"github.com/cjlapao/common-go/identity/authorization_context"
	"github.com/cjlapao/common-go/identity/constants"
	"github.com/cjlapao/common-go/identity/jwt"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/service_provider"
	"github.com/gorilla/mux"
)

// TokenAuthorizationMiddlewareAdapter validates a Authorization Bearer during a rest api call
// It can take an array of roles and claims to further validate the token in a more granular
// view, it also can take an OR option in both if the role or claim are coma separated.
// For example a claim like "_read,_write" will be valid if the user either has a _read claim
// or a _write claim making them both valid
func TokenAuthorizationMiddlewareAdapter(roles []string, claims []string) controllers.Adapter {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := execution_context.Get()
			if ctx.UserDatabaseAdapter == nil {
				w.WriteHeader(http.StatusUnauthorized)
				identity.ErrNoContext.Log()
				json.NewEncoder(w).Encode(identity.ErrNoContext)
				return
			}
			vars := mux.Vars(r)
			tenantId := vars["tenantId"]
			var userToken *models.UserToken
			var dbUser *models.User
			var validateError error
			var err error
			var userClaims = make([]string, 0)
			var userRoles = make([]string, 0)
			var isSuperUser = false

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
				logger.Error("Error validating token, %v", validateError.Error())
			}

			// Validating userToken against the keys
			if authorized {
				var validateUserTokenError error
				if ctx.Authorization.Options.KeyVaultEnabled {
					userToken, validateUserTokenError = jwt.ValidateUserToken(jwt_token, ctx.Authorization.Scope)
				} else if ctx.Authorization.Options.PublicKey != "" {
					userToken, validateUserTokenError = jwt.ValidateUserToken(jwt_token, ctx.Authorization.Scope)
				} else {
					validateUserTokenError = errors.New("No public or private key found to validate token")
				}

				if validateUserTokenError != nil {
					authorized = false
					validateError = errors.New("bearer token is not valid, " + validateUserTokenError.Error())
					logger.Error("Error validating token, %v", validateError.Error())
				}
			}

			// Making sure the user token does contain the necessary fields
			if userToken == nil || userToken.User == "" {
				authorized = false
				validateError = errors.New("bearer token has invalid user")
				logger.Error("Error validating token, %v", validateError.Error())
			}

			// Checking if the user is a supper user, if so we will not check any roles or claims as he will have access to it
			if len(roles) > 0 {
				for _, role := range roles {
					if role == constants.SuperUser {
						logger.Info("Super User %v was found, authorizing", userToken.User)
						authorized = true
						isSuperUser = true
						break
					}
				}
			}

			// To gain speed we will only get the db user if there is any role or claim to validate
			// otherwise we don't need anything else to validate it
			if (len(roles) > 0 || len(claims) > 0) && !isSuperUser {
				// Getting the user from the database to validate roles and claims
				if authorized {
					dbUser = ctx.UserDatabaseAdapter.GetUserByEmail(userToken.User)
					if dbUser == nil || dbUser.ID == "" {
						authorized = false
						validateError = errors.New("bearer token user was not found in database, potentially revoked, " + userToken.User)
						logger.Error("Error validating token, %v", validateError.Error())
					}
				}

				if authorized {
					userRoles, userClaims = getUserRolesAndClaims(dbUser)
				}

				// Validating user roles
				if authorized && len(roles) > 0 && len(userRoles) > 0 {
					err = validateUserRoles(userRoles, roles)
					if err != nil {
						authorized = false
						validateError = errors.New("bearer token user does not contain one or more roles required by the context, " + err.Error())
						logger.Error("Error validating token, %v", validateError.Error())
					}
				}

				// Validating user claims
				if authorized && len(claims) > 0 && len(userClaims) > 0 {
					err = validateUserClaims(userClaims, claims)
					if err != nil {
						authorized = false
						validateError = errors.New("bearer token user does not contain one or more claims required by the context, " + err.Error())
						logger.Error("Error validating token, %v", validateError.Error())
					}
				}
			}

			if authorized && userToken != nil && userToken.ID != "" {
				user := authorization_context.NewUserContext()
				user.ID = userToken.ID
				user.Email = userToken.User
				user.Audiences = userToken.Audiences
				user.Issuer = userToken.Issuer
				user.ValidatedClaims = claims
				user.Roles = userToken.Roles

				baseUrl := service_provider.Get().GetBaseUrl(r)
				ctx.Authorization = authorization_context.NewFromUser(user)
				ctx.Authorization.Issuer = baseUrl + "/auth/" + tenantId
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
					logger.Error("User "+userToken.User+" failed to authorize, %v", response.ErrorDescription)
				} else {
					logger.Error("Request failed to authorize, %v", response.ErrorDescription)

				}

				logger.Info("Authorization Layer Finished")
				return
			}

			logger.Info("Authorization Layer Finished")
			next.ServeHTTP(w, r)
		})
	}
}

// EndAuthorizationMiddlewareAdapter This cleans the context of any previous users
// token left in memory and rereading all of the default options for the next request
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

func getUserRolesAndClaims(user *models.User) (roles []string, claims []string) {
	roles = make([]string, 0)
	claims = make([]string, 0)

	if user == nil || user.ID == "" {
		return roles, claims
	}

	if !user.IsValid() {
		return roles, claims
	}

	for _, role := range user.Roles {
		roles = append(roles, role.ID)
	}

	for _, claim := range user.Claims {
		claims = append(claims, claim.ID)
	}

	return roles, claims
}

func validateUserRoles(roles []string, requiredRoles []string) error {
	var validateError error

	// Getting user roles and claims
	validatedRoles := make(map[string]bool)

	if len(requiredRoles) > 0 {
		if len(roles) == 0 {
			validateError = errors.New("user does not contain any roles")
			logger.Error(validateError.Error())
			return validateError
		} else {
			for _, requiredRole := range requiredRoles {
				foundRole := false
				for _, role := range roles {
					requiredRoleArr := strings.Split(requiredRole, ",")
					if len(requiredRoleArr) == 1 {
						if strings.EqualFold(requiredRole, role) {
							foundRole = true
							break
						}
					} else if len(requiredRoleArr) > 1 {
						for _, splitRequiredRole := range requiredRoleArr {
							if strings.EqualFold(splitRequiredRole, role) {
								foundRole = true
								break
							}
						}
						if foundRole {
							break
						}
					}
				}

				validatedRoles[requiredRole] = foundRole
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

func validateUserClaims(claims []string, requiredClaims []string) error {
	var validateError error

	// Getting user claims and claims
	validatedClaims := make(map[string]bool)

	if len(requiredClaims) > 0 {
		if len(claims) == 0 {
			validateError = errors.New("user does not contain any claims")
			logger.Error(validateError.Error())
			return validateError
		} else {
			for _, requiredClaim := range requiredClaims {
				foundClaim := false
				for _, claim := range claims {
					requiredClaimArr := strings.Split(requiredClaim, ",")
					if len(requiredClaimArr) == 1 {
						if strings.EqualFold(requiredClaim, claim) {
							foundClaim = true
							break
						}
					} else if len(requiredClaimArr) > 1 {
						for _, splitRequiredClaim := range requiredClaimArr {
							if strings.EqualFold(splitRequiredClaim, claim) {
								foundClaim = true
								break
							}
						}
						if foundClaim {
							break
						}
					}
				}

				validatedClaims[requiredClaim] = foundClaim
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
