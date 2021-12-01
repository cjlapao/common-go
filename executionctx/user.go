package executionctx

type UserCtx struct {
	Name      string
	Audiences []string
	Issuer    string
}

func NewUserContext(user *UserCtx) {
	ctx := GetContext()
	ctx.User = user
}

func ClearUserContext() {
	ctx := GetContext()
	ctx.User = nil
}
