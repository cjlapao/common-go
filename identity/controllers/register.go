package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/cjlapao/common-go/controllers"
	"github.com/cjlapao/common-go/helper/http_helper"
	"github.com/cjlapao/common-go/identity/constants"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/security"
)

// Register Create an user in the tenant
func (c *AuthorizationControllers) Register() controllers.Controller {
	return func(w http.ResponseWriter, r *http.Request) {
		var registerRequest models.OAuthRegisterRequest

		http_helper.MapRequestBody(r, &registerRequest)

		user := models.NewUser()
		user.Username = registerRequest.Username
		user.Email = registerRequest.Email
		user.FirstName = registerRequest.FirstName
		user.LastName = registerRequest.LastName
		user.DisplayName = user.FirstName + " " + user.LastName
		user.Password = security.SHA256Encode(registerRequest.Password)
		user.InvalidAttempts = 0
		user.EmailVerified = false

		if registerRequest.Claims != nil && len(registerRequest.Claims) > 0 {
			for _, claim := range registerRequest.Claims {
				user.Claims = append(user.Claims, models.NewUserClaim(claim, claim))
			}
		} else {
			user.Claims = append(user.Claims, constants.ReadClaim)
		}

		if registerRequest.Roles != nil && len(registerRequest.Roles) > 0 {
			for _, role := range registerRequest.Roles {
				user.Roles = append(user.Roles, models.NewUserRole(role, role))
			}
		} else {
			user.Roles = append(user.Roles, constants.RegularUserRole)
		}

		if !user.IsValid() {
			w.WriteHeader(http.StatusUnauthorized)
			ErrInvalidUser.Log()
			json.NewEncoder(w).Encode(ErrInvalidUser)
			return
		}

		dbUser := c.Context.UserDatabaseAdapter.GetUserByEmail(user.Email)
		if dbUser.Email != "" {
			w.WriteHeader(http.StatusBadRequest)
			ErrUserAlreadyExists.Log()
			json.NewEncoder(w).Encode(ErrUserAlreadyExists)
			return
		}

		c.Context.UserDatabaseAdapter.UpsertUser(*user)

		json.NewEncoder(w).Encode(user)
	}
}
