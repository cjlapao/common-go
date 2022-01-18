package validators

type validField struct {
	characters map[string]bool
}

func New(characters string) *validField {

	var obj = validField{make(map[string]bool, len(characters))}

	// build map from runes
	for _, val := range characters {
		obj.characters[string(val)] = true
	}

	return &obj
}

func (v *validField) ValidateField(value string) bool {
	return v.characters[value]
}
