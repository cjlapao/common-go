package controllers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity"
	"github.com/cjlapao/common-go/identity/authorization_context"
	"github.com/cjlapao/common-go/identity/identity_database_adapter"
	"github.com/cjlapao/common-go/identity/jwt"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security"
)

type AuthorizationControllers interface {
	Context() authorization_context.AuthorizationContext
	Login() controllers.Controller
	Validate() controllers.Controller
}

type DefaultAuthorizationControllers struct {
	UserAdapter identity_database_adapter.UserDatabaseAdapter
}

func NewDefaultAuthorizationControllers() *DefaultAuthorizationControllers {
	context := identity_database_adapter.NewMemoryUserAdapter()
	controllers := DefaultAuthorizationControllers{
		UserAdapter: context,
	}
	return &controllers
}

func NewAuthorizationControllers(context identity_database_adapter.UserDatabaseAdapter) *DefaultAuthorizationControllers {
	controllers := DefaultAuthorizationControllers{
		UserAdapter: context,
	}
	return &controllers
}

// Login Generate a token for a valid user
func (c *DefaultAuthorizationControllers) Login() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		loginRequestUser := models.LoginRequest{}
		json.Unmarshal(reqBody, &loginRequestUser)

		if c.UserAdapter == nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("There was an error during login, context is null")
			return
		}
		user := c.UserAdapter.GetUserByEmail(loginRequestUser.Username)

		if user.ID == "" {
			w.WriteHeader(http.StatusUnauthorized)
			if user.DisplayName != "" {
				logger.Error("There was an error during login, user %v was not found", user.DisplayName)
			} else if len(reqBody) == 0 {
				logger.Error("There was an error during login, body was empty")
			} else {
				logger.Error("There was an error during login, unknown error")
			}
			return
		}

		password := security.SHA256Encode(loginRequestUser.Password)

		if password != user.Password {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("There was an error during loggin user %v, password is incorrect", user.Username)
			return
		}

		token, err := jwt.GenerateDefaultUserToken(*user)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("There was an error during login, error generating the key")
			return
		}

		c.UserAdapter.UpdateUserRefreshToken(user.ID, token.RefreshToken)

		response := models.LoginResponse{
			AccessToken:  token.Token,
			RefreshToken: token.RefreshToken,
			ExpiresIn:    token.ExpiresAt,
			CreatedAt:    time.Now(),
		}

		logger.Success("User %v was logged in successfully", user.Username)

		json.NewEncoder(w).Encode(response)
	}
}

// Validate Validate a token for a valid user
func (c *DefaultAuthorizationControllers) Validate() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		token, valid := http_helper.GetAuthorizationToken(r.Header)

		if !valid {
			response := models.LoginErrorResponse{
				Error:            "Token Not Found",
				ErrorDescription: "The JWT token was not found or the header was malformed, please check your authorization header",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			logger.Error("There was an error validating token")
			return
		}

		if !security.ValidateToken(token) {
			response := models.LoginErrorResponse{
				Error:            "Invalid Token",
				ErrorDescription: "The JWT token is not valid",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			logger.Error("There was an error validating token")
			return
		}

		response := models.LoginResponse{
			AccessToken: token,
		}

		logger.Success("Token was validated successfully")
		json.NewEncoder(w).Encode(response)
	}
}

// Validate Validate a token for a valid user
func (c *DefaultAuthorizationControllers) Introspection() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		token, valid := http_helper.GetAuthorizationToken(r.Header)

		if !valid {
			response := models.LoginErrorResponse{
				Error:            models.TokenNotFoundError,
				ErrorDescription: "The JWT token was not found or the header was malformed, please check your authorization header",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			logger.Error("There was an error validating token")
			return
		}

		isTokenValid, err := jwt.ValidateUserToken(token, identity.ApplicationTokenScope)

		if !isTokenValid {
			response := models.LoginErrorResponse{
				Error:            models.InvalidTokenError,
				ErrorDescription: err.Error(),
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			logger.Error("There was an error validating token")
			return
		}

		response := models.LoginResponse{
			AccessToken: token,
		}

		logger.Success("Token was validated successfully")
		json.NewEncoder(w).Encode(response)
	}
}
