package constants

import "github.com/cjlapao/common-go/identity/models"

const (
	SuperUserClaim        = "_su"
	CanUpsertClaim        = "_readwrite"
	CanReadClaim          = "_read"
	CanWriteClaim         = "_write"
	CanRemoveClaim        = "_remove"
	CanUpsertUserClaim    = "_readwrite.user"
	CanReadUserClaim      = "_read.user"
	CanReadWriteUserClaim = "_readwrite.user"
	CanRemoveUserClaim    = "_remove.user"
)

var ReadClaim = models.UserClaim{
	ID:   CanReadClaim,
	Name: "Can Read",
}

var ReadWriteClaim = models.UserClaim{
	ID:   CanUpsertClaim,
	Name: "Can Read/Write",
}

var RemoveClaim = models.UserClaim{
	ID:   CanRemoveClaim,
	Name: "Can Remove",
}
