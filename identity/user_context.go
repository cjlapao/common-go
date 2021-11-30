package identity

import "strings"

type UserContext interface {
	GetUserById(id string) *User
	GetUserByEmail(email string) *User
	UpsertUser(user User)
}

type DefaultUserContext struct {
	Users []User
}

func NewDefaultUserContext() *DefaultUserContext {
	context := DefaultUserContext{}
	context.Users = GetDefaultUsers()

	return &context
}

func (c *DefaultUserContext) GetUserById(id string) *User {
	users := GetDefaultUsers()
	var user User
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

func (c *DefaultUserContext) GetUserByEmail(email string) *User {
	users := GetDefaultUsers()
	var user User
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

func (c *DefaultUserContext) UpsertUser(user User) {
	c.Users = append(c.Users, user)
}
