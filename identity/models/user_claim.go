package models

import (
	"github.com/cjlapao/common-go/identity/constants"
	"github.com/google/uuid"
)

type UserClaim struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"claimName" bson:"claimName"`
}

func NewUserClaim(name string) UserClaim {
	return UserClaim{
		ID:   uuid.NewString(),
		Name: name,
	}
}

func (uc UserClaim) IsValid() bool {
	_, err := uuid.Parse(uc.ID)
	if err != nil {
		return false
	}
	if uc.Name == "" {
		return false
	}

	return true
}

var ReadClaim = UserClaim{
	ID:   constants.CanReadClaim,
	Name: "Can View",
}

var ReadWriteClaim = UserClaim{
	ID:   constants.CanUpsertClaim,
	Name: "Can Edit/Change",
}

var RemoveClaim = UserClaim{
	ID:   constants.CanUpsertClaim,
	Name: "Can Remove",
}
