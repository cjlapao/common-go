package identity

import (
	"fmt"

	"github.com/cjlapao/common-go/database/mongodb"
	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/security"
	"github.com/google/uuid"
)

func GetDefaultUsers() []User {
	config := executionctx.GetConfigService()
	users := make([]User, 0)
	var adminUser User
	var demoUser User

	adminUsername := config.GetString("ADMIN_USERNAME")
	adminPassword := config.GetString("ADMIN_PASSWORD")

	if helper.IsNilOrEmpty(adminUsername) {
		adminUser = User{
			ID:       uuid.NewString(),
			Email:    "admin@localhost",
			Username: "admin",
			Password: "a075d17f3d453073853f813838c15b8023b8c487038436354fe599c3942e1f95",
		}
	} else {
		adminUser = User{
			ID:       uuid.NewString(),
			Email:    fmt.Sprintf("%v@localhost", adminUsername),
			Username: adminUsername,
		}
	}

	if helper.IsNilOrEmpty(adminPassword) {
		security.SHA256Encode("p@ssw0rd")
	} else {
		adminUser.Password = security.SHA256Encode(adminPassword)
	}
	roles := []UserRole{AdminRole, RegularUserRole}
	adminUser.Roles = append(adminUser.Roles, roles...)

	demoUsername := config.GetString("DEMO_USERNAME")
	demoPassword := config.GetString("DEMO_PASSWORD")

	if helper.IsNilOrEmpty(adminUsername) {
		demoUser = User{
			ID:       uuid.NewString(),
			Email:    "demo@localhost",
			Username: "demo",
			Password: "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea",
		}
	} else {
		demoUser = User{
			ID:       uuid.NewString(),
			Email:    fmt.Sprintf("%v@localhost", demoUsername),
			Username: adminUsername,
		}
	}

	if helper.IsNilOrEmpty(demoPassword) {
		security.SHA256Encode("demo")
	} else {
		demoUser.Password = security.SHA256Encode(demoPassword)
	}

	demoUser.Roles = append(demoUser.Roles, RegularUserRole)

	users = append(users, adminUser)
	users = append(users, demoUser)

	return users
}

func Seed(factory *mongodb.MongoFactory, databaseName string) {
	SeedUsers(factory, databaseName)
}

func SeedUsers(factory *mongodb.MongoFactory, databaseName string) {
	repo := mongodb.NewRepository(factory, databaseName, IdentityUsersCollection)
	users := GetDefaultUsers()
	for _, user := range users {
		filter := make(map[string]interface{})
		filter["email"] = user.Email
		repo.UpsertOne(filter, user)
	}
}
