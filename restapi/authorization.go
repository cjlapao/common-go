package restapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/security"
)

type AuthorizationContext struct {
	User string
}

type AuhtorizationController struct{}

// Login Generate a token for a valid user
func (l *HttpListener) Login(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	loginRequestUser := LoginRequest{}
	json.Unmarshal(reqBody, &loginRequestUser)
	users := GetDefaultUsers()
	var user User
	found := false
	for _, usr := range users {
		if strings.ToLower(loginRequestUser.Username) == strings.ToLower(usr.Username) {
			user = usr
			found = true
			break
		}
	}

	if !found {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	password := security.SHA256Encode(loginRequestUser.Password)

	if password != user.Password {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, expires := security.GenerateUserToken(user.Email)
	response := LoginResponse{
		AccessToken: string(token),
		Expiring:    expires,
	}

	json.NewEncoder(w).Encode(response)
}

// Validate Validate a token for a valid user
func (l *HttpListener) Validate(w http.ResponseWriter, r *http.Request) {
	token, valid := security.GetAuthorizationToken(r.Header)

	if !valid {
		response := LoginErrorResponse{
			Code:    "404",
			Error:   "Token Not Found",
			Message: "The JWT token was not found or the header was malformed, please check your authorization header",
		}

		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
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
		return
	}

	response := LoginResponse{
		AccessToken: token,
	}

	json.NewEncoder(w).Encode(response)
}

func GetDefaultUsers() []User {
	config := executionctx.GetConfigService()
	users := make([]User, 0)
	var adminUser User
	var demoUser User

	adminUsername := config.GetString("ADMIN_USERNAME")
	adminPassword := config.GetString("ADMIN_PASSWORD")

	if helper.IsNilOrEmpty(adminUsername) {
		adminUser = User{
			Email:    "admin@localhost",
			Username: "admin",
			Password: "a075d17f3d453073853f813838c15b8023b8c487038436354fe599c3942e1f95",
		}
	} else {
		adminUser = User{
			Email:    fmt.Sprintf("%v@localhost", adminUsername),
			Username: adminUsername,
		}
	}

	if helper.IsNilOrEmpty(adminPassword) {
		security.SHA256Encode("p@ssw0rd")
	} else {
		adminUser.Password = security.SHA256Encode(adminPassword)
	}

	demoUsername := config.GetString("DEMO_USERNAME")
	demoPassword := config.GetString("DEMO_PASSWORD")

	if helper.IsNilOrEmpty(adminUsername) {
		demoUser = User{
			Email:    "demo@localhost",
			Username: "demo",
			Password: "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea",
		}
	} else {
		demoUser = User{
			Email:    fmt.Sprintf("%v@localhost", demoUsername),
			Username: adminUsername,
		}
	}

	if helper.IsNilOrEmpty(demoPassword) {
		security.SHA256Encode("demo")
	} else {
		demoUser.Password = security.SHA256Encode(demoPassword)
	}

	users = append(users, adminUser)
	users = append(users, demoUser)

	return users
}
