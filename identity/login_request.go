package identity

// LoginRequest entity
type LoginRequest struct {
	Username string `json:"username" bson:"username" yaml:"username"`
	Password string `json:"password" bson:"password" yaml:"username"`
}
