package mongodb

type FilterElement struct {
	Key   string
	Value interface{}
}

func (e FilterElement) Encode() map[string]interface{} {
	mapped := make(map[string]interface{})
	mapped[e.Key] = e.Value
	return mapped
}
