package identity

import (
	"fmt"

	"github.com/cjlapao/common-go/executionctx"
	"github.com/cjlapao/common-go/helper"
	"github.com/cjlapao/common-go/security"
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
			Email:    "admin@localhost",
			Username: "admin",
			Password: "a075d17f3d453073853f813838c15b8023b8c487038436354fe599c3942e1f95",
		}
	} else {
		adminUser = User{
			Email:    fmt.Sprintf("%v@localhost", adminUsername),
			Username: adminUsername,
		}
	}

	if helper.IsNilOrEmpty(adminPassword) {
		security.SHA256Encode("p@ssw0rd")
	} else {
		adminUser.Password = security.SHA256Encode(adminPassword)
	}

	demoUsername := config.GetString("DEMO_USERNAME")
	demoPassword := config.GetString("DEMO_PASSWORD")

	if helper.IsNilOrEmpty(adminUsername) {
		demoUser = User{
			Email:    "demo@localhost",
			Username: "demo",
			Password: "2a97516c354b68848cdbd8f54a226a0a55b21ed138e207ad6c5cbb9c00aa5aea",
		}
	} else {
		demoUser = User{
			Email:    fmt.Sprintf("%v@localhost", demoUsername),
			Username: adminUsername,
		}
	}

	if helper.IsNilOrEmpty(demoPassword) {
		security.SHA256Encode("demo")
	} else {
		demoUser.Password = security.SHA256Encode(demoPassword)
	}

	users = append(users, adminUser)
	users = append(users, demoUser)

	return users
}
