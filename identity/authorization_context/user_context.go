package authorization_context

import (
	"github.com/cjlapao/common-go/constants"
	"github.com/cjlapao/common-go/cryptorand"
)

type UserContext struct {
	ID              string
	Username        string
	Email           string
	DisplayName     string
	Tenant          string
	Audiences       []string
	Issuer          string
	ValidatedClaims []string
	Roles           []string
}

func NewUserContext() *UserContext {
	user := UserContext{
		ID:              cryptorand.GenerateRandomString(constants.ID_SIZE),
		ValidatedClaims: make([]string, 0),
		Roles:           make([]string, 0),
	}

	return &user
}
