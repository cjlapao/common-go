package identity_database_adapter

import (
	"strings"

	"github.com/cjlapao/common-go/identity"
	"github.com/cjlapao/common-go/identity/models"
)

type MemoryUserAdapter struct {
	Users []models.User
}

func NewMemoryUserAdapter() *MemoryUserAdapter {
	context := MemoryUserAdapter{}
	context.Users = identity.GetDefaultUsers()

	return &context
}

func (c *MemoryUserAdapter) GetUserById(id string) *models.User {
	users := identity.GetDefaultUsers()
	var user models.User
	found := false
	for _, usr := range users {
		if strings.EqualFold(id, usr.ID) {
			user = usr
			found = true
			break
		}
	}

	if found {
		return &user
	}
	return nil
}

func (c *MemoryUserAdapter) GetUserByEmail(email string) *models.User {
	users := identity.GetDefaultUsers()
	var user models.User
	found := false
	for _, usr := range users {
		if strings.EqualFold(email, usr.Email) {
			user = usr
			found = true
			break
		}
	}

	if found {
		return &user
	}
	return nil
}

func (c *MemoryUserAdapter) UpsertUser(user models.User) {
	c.Users = append(c.Users, user)
}
