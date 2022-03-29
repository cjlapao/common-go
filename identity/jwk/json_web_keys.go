package jwk

type JsonWebKeys struct {
	Keys []*JsonWebKey `json:"keys"`
}

func New() *JsonWebKeys {
	keys := JsonWebKeys{}
	keys.Keys = make([]*JsonWebKey, 0)

	return &keys
}

func (jk *JsonWebKeys) Add(id string, privateKey interface{}) {
	key := NewKeyWithId(id, privateKey)
	jk.Keys = append(jk.Keys, key)
}
