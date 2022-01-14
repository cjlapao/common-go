package constants

import "github.com/cjlapao/common-go/identity/models"

const (
	SuperUser   = "_su"
	Admin       = "_admin"
	RegularUser = "_user"
)

var SuRole = models.UserRole{
	ID:   SuperUser,
	Name: "Super Administrator",
}

var AdminRole = models.UserRole{
	ID:   Admin,
	Name: "Administrator",
}

var RegularUserRole = models.UserRole{
	ID:   RegularUser,
	Name: "User",
}
