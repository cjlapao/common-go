package interfaces

import (
	"github.com/cjlapao/common-go/identity/models"
)

type UserDatabaseAdapter interface {
	GetUserById(id string) *models.User
	GetUserByEmail(email string) *models.User
	GetUserByUsername(username string) *models.User
	UpsertUser(user models.User)
	RemoveUser(id string) bool
	GetUserRefreshToken(id string) *string
	UpdateUserRefreshToken(id string, token string) bool
	GetUserEmailVerifyToken(id string) *string
	UpdateUserEmailVerifyToken(id string, token string) bool
}
