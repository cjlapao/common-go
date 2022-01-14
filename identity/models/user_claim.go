package models

type UserClaim struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"claimName" bson:"claimName"`
}

func NewUserClaim(id string, name string) UserClaim {
	return UserClaim{
		ID:   id,
		Name: name,
	}
}

func (uc UserClaim) IsValid() bool {
	return uc.ID != "" && uc.Name != ""
}
