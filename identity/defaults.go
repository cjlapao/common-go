package identity

import (
	"fmt"

	"github.com/cjlapao/common-go/configuration"
	"github.com/cjlapao/common-go/database/mongodb"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/identity/models"
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/security"
)

var logger = log.Get()

func GetDefaultUsers() []models.User {
	config := configuration.Get()
	users := make([]models.User, 0)
	var adminUser models.User
	var demoUser models.User

	adminUsername := config.GetString("ADMIN_USERNAME")
	adminPassword := config.GetString("ADMIN_PASSWORD")

	if helper.IsNilOrEmpty(adminUsername) {
		adminUser = models.User{
			ID:        "592D8E5C-6F5D-40A0-9348-80131B083715",
			Email:     "admin@localhost",
			FirstName: "Administrator",
			LastName:  "User",
			Username:  "admin",
			Password:  "a075d17f3d453073853f813838c15b8023b8c487038436354fe599c3942e1f95",
		}
	} else {
		adminUser = models.User{
			ID:        "592D8E5C-6F5D-40A0-9348-80131B083715",
			Email:     fmt.Sprintf("%v@localhost", adminUsername),
			FirstName: "Administrator",
			LastName:  "User",
			Username:  adminUsername,
		}
	}

	if helper.IsNilOrEmpty(adminPassword) {
		security.SHA256Encode("p@ssw0rd")
	} else {
		adminUser.Password = security.SHA256Encode(adminPassword)
	}
	roles := []models.UserRole{models.AdminRole, models.RegularUserRole}
	adminUser.Roles = append(adminUser.Roles, roles...)

	demoUsername := config.GetString("DEMO_USERNAME")
	demoPassword := config.GetString("DEMO_PASSWORD")

	if helper.IsNilOrEmpty(demoUsername) {
		demoUser = models.User{
			ID:        "C54C2A9B-CA73-4188-875A-F26026A38B58",
			Email:     "demo@localhost",
			FirstName: "Demo",
			LastName:  "User",
			Username:  "demo",
			Password:  "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea",
		}
	} else {
		demoUser = models.User{
			ID:        "C54C2A9B-CA73-4188-875A-F26026A38B58",
			Email:     fmt.Sprintf("%v@localhost", demoUsername),
			FirstName: "Administrator",
			LastName:  "User",
			Username:  demoUsername,
		}
	}

	if helper.IsNilOrEmpty(demoPassword) {
		security.SHA256Encode("demo")
	} else {
		demoUser.Password = security.SHA256Encode(demoPassword)
	}

	demoUser.Roles = append(demoUser.Roles, models.RegularUserRole)

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
		model := mongodb.NewUpdateOneBuilder().FilterBy("email", user.Email).Encode(user).Build()
		repo.UpsertOne(model)
	}
}
