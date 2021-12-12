package identity

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/security"
)

type AuthorizationControlles interface {
	Context() UserContext
	Login() controllers.Controller
	Validate() controllers.Controller
}

type AuthorizationContext struct {
	User User
}

type DefaultAuthorizationControllers struct {
	Context UserContext
}

func NewDefaultAuthorizationControllers() *DefaultAuthorizationControllers {
	context := NewDefaultUserContext()
	controllers := DefaultAuthorizationControllers{
		Context: context,
	}
	return &controllers
}

func NewAuthorizationControllers(context UserContext) *DefaultAuthorizationControllers {
	controllers := DefaultAuthorizationControllers{
		Context: context,
	}
	return &controllers
}

// Login Generate a token for a valid user
func (c *DefaultAuthorizationControllers) Login() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBody, _ := ioutil.ReadAll(r.Body)
		loginRequestUser := LoginRequest{}
		json.Unmarshal(reqBody, &loginRequestUser)

		if c.Context == nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("There was an error during loggin, context is null")
			return
		}
		user := c.Context.GetUserByEmail(loginRequestUser.Username)

		if user == nil {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("There was an error during loggin, user %v was not found", user.Username)
			return
		}

		password := security.SHA256Encode(loginRequestUser.Password)

		if password != user.Password {
			w.WriteHeader(http.StatusUnauthorized)
			logger.Error("There was an error during loggin user %v, password is incorrect", user.Username)
			return
		}

		token, expires := security.GenerateUserToken(user.Email)
		response := LoginResponse{
			AccessToken: string(token),
			Expiring:    expires,
		}
		logger.Success("User %v was logged in successfully", user.Username)

		json.NewEncoder(w).Encode(response)
	}
}

// Validate Validate a token for a valid user
func (c *DefaultAuthorizationControllers) Validate() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		token, valid := security.GetAuthorizationToken(r.Header)

		if !valid {
			response := LoginErrorResponse{
				Code:    "404",
				Error:   "Token Not Found",
				Message: "The JWT token was not found or the header was malformed, please check your authorization header",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			logger.Error("There was an error validating token")
			return
		}

		if !security.ValidateToken(token) {
			response := LoginErrorResponse{
				Code:    "401",
				Error:   "Invalid Token",
				Message: "The JWT token is not valid",
			}

			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(response)
			logger.Error("There was an error validating token")
			return
		}

		response := LoginResponse{
			AccessToken: token,
		}

		logger.Success("Token was validated successfully")
		json.NewEncoder(w).Encode(response)
	}
}
