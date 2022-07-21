package authorization_context

import (
	cryptorand "github.com/cjlapao/common-go-cryptorand"
	"github.com/cjlapao/common-go/constants"
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
		ID:              cryptorand.GetRandomString(constants.ID_SIZE),
		ValidatedClaims: make([]string, 0),
		Roles:           make([]string, 0),
	}

	return &user
}
