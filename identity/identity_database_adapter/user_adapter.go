package identity_database_adapter

import (
	"github.com/cjlapao/common-go/identity/models"
)

type UserDatabaseAdapter interface {
	GetUserById(id string) *models.User
	GetUserByEmail(email string) *models.User
	UpsertUser(user models.User)
}
