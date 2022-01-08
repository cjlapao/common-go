package database

import (
	"strings"

	"github.com/cjlapao/common-go/identity"
	"github.com/cjlapao/common-go/identity/models"
)

type MemoryUserContextAdapter struct {
	Users []models.User
}

func NewMemoryUserAdapter() *MemoryUserContextAdapter {
	context := MemoryUserContextAdapter{}
	context.Users = identity.GetDefaultUsers()

	return &context
}

func (c *MemoryUserContextAdapter) GetUserById(id string) *models.User {
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

func (c *MemoryUserContextAdapter) GetUserByEmail(email string) *models.User {
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

func (c *MemoryUserContextAdapter) UpsertUser(user models.User) {
	c.Users = append(c.Users, user)
}

func (c *MemoryUserContextAdapter) GetUserRefreshToken(id string) *string {
	user := c.GetUserById(id)
	token := ""
	if user != nil {
		token = user.RefreshToken
	}

	return &token
}

func (c *MemoryUserContextAdapter) UpdateUserRefreshToken(id string, token string) {
	user := c.GetUserById(id)
	if user != nil {
		user.RefreshToken = token
	}
}

func (c *MemoryUserContextAdapter) GetUserEmailVerifyToken(id string) *string {
	user := c.GetUserById(id)
	token := ""
	if user != nil {
		token = user.EmailVerifyToken
	}

	return &token
}

func (c *MemoryUserContextAdapter) UpdateUserEmailVerifyToken(id string, token string) {
	user := c.GetUserById(id)
	if user != nil {
		user.EmailVerifyToken = token
	}
}
