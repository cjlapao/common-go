package controllers

import (
	"github.com/cjlapao/common-go/identity/database"
	"github.com/cjlapao/common-go/identity/interfaces"
	"github.com/cjlapao/common-go/log"
)

var logger = log.Get()

type AuthorizationControllers struct {
	UserAdapter interfaces.UserDatabaseAdapter
}

func NewDefaultAuthorizationControllers() *AuthorizationControllers {
	context := database.NewMemoryUserAdapter()
	controllers := AuthorizationControllers{
		UserAdapter: context,
	}
	return &controllers
}

func NewAuthorizationControllers(context interfaces.UserDatabaseAdapter) *AuthorizationControllers {
	controllers := AuthorizationControllers{
		UserAdapter: context,
	}
	return &controllers
}
