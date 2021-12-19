package authorization_context

import "github.com/google/uuid"

type ContextUser struct {
	ID          string
	Username    string
	Email       string
	DisplayName string
	Tenant      string
	Audiences   []string
	Issuer      string
	Claims      []string
}

func NewContextUser() *ContextUser {
	user := ContextUser{
		ID: uuid.NewString(),
	}

	return &user
}
