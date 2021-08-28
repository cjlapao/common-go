package controller

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"

// 	"github.com/cjlapao/common-go/executionctx"
// 	"github.com/cjlapao/common-go/log"
// 	"github.com/cjlapao/common-go/security"
// )

// // Login Generate a token for a valid user
// func (c *Controller) Probe(w http.ResponseWriter, r *http.Request) {
// 	logger := log.Get()
// 	config := executionctx.GetConfigService()

// 	port := config.Get("port")
// 	fmt.Print(port)
// 	apiPort := config.Get("LOADTEST_DEBUG")
// 	fmt.Print(apiPort)

// 	logger.Debug("Login Endpoint Hit")

// 	reqBody, _ := ioutil.ReadAll(r.Body)
// 	loginRequest := LoginRequest{}
// 	json.Unmarshal(reqBody, &loginRequest)

// 	user := User{
// 		Email:    "admin@localhost",
// 		Username: "admin",
// 		Password: "a075d17f3d453073853f813838c15b8023b8c487038436354fe599c3942e1f95",
// 	}
// 	if len(user.Email) == 0 {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	password := security.SHA256Encode(loginRequest.Password)

// 	if password != user.Password {
// 		w.WriteHeader(http.StatusUnauthorized)
// 		return
// 	}

// 	token, expires := security.GenerateUserToken(user.Email)
// 	response := LoginResponse{
// 		AccessToken: string(token),
// 		Expiring:    expires,
// 	}

// 	json.NewEncoder(w).Encode(response)
// }
