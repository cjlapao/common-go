package controllers

import (
	"github.com/cjlapao/common-go/execution_context"
	"github.com/cjlapao/common-go/identity/database"
	"github.com/cjlapao/common-go/identity/interfaces"
	"github.com/cjlapao/common-go/log"
)

// AuthorizationControllers
type AuthorizationControllers struct {
	Logger  *log.Logger
	Context *execution_context.Context
}

func NewDefaultAuthorizationControllers() *AuthorizationControllers {
	ctx := execution_context.Get()
	context := database.NewMemoryUserAdapter()
	ctx.UserDatabaseAdapter = context
	controllers := AuthorizationControllers{
		Logger:  log.Get(),
		Context: ctx,
	}

	return &controllers
}

func NewAuthorizationControllers(context interfaces.UserDatabaseAdapter) *AuthorizationControllers {
	ctx := execution_context.Get()
	ctx.UserDatabaseAdapter = context
	controllers := AuthorizationControllers{
		Logger:  log.Get(),
		Context: ctx,
	}

	return &controllers
}
