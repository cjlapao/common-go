package cryptorand

import "github.com/cjlapao/common-go/constants"

func GenerateAlphaNumericRandomString(size int) string {
	return randomString(size, true, true, true)
}

func GenerateRandomString(size int) string {
	return randomString(size, true, true, false)
}

func GenerateUpperCaseRandomString(size int) string {
	return randomString(size, false, true, false)
}

func GenerateLowerCaseRandomString(size int) string {
	return randomString(size, true, false, false)
}

func GenerateNumericRandomString(size int) string {
	return randomString(size, false, false, true)
}

func randomString(size int, includeLowerCase bool, includeUpperCase bool, includeNumeric bool) string {
	source := make([]string, 0)
	result := ""

	if includeLowerCase {
		source = append(source, constants.LowerCaseAlphaCharacters()...)
	}

	if includeUpperCase {
		source = append(source, constants.UpperCaseAlphaCharacters()...)
	}

	if includeNumeric {
		source = append(source, constants.NumericCharacters()...)
	}
	if len(source) > 0 {
		random := Rand()
		if size > 0 {
			for i := 0; i < size; i++ {
				idx := random.Intn(len(source))
				result += source[idx]
			}
		}
	}

	return result
}
