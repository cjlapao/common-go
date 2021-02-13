package helper

import (
	"encoding/hex"
	"math/rand"
	"time"
)

// RandomHex generates a random hex string
func RandomHex(size int) string {
	randString := RandomString(size / 2)

	randStringBytes := []byte(randString)

	return hex.EncodeToString(randStringBytes)
}

// RandomHexWithPrefix generates a random hex string
func RandomHexWithPrefix(size int) string {
	return "0x" + RandomHex(size)
}

// RandomString generates a random string
func RandomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	var result []rune
	source := []rune(AlphaNumeric)
	for i := 0; i < size; i++ {
		result = append(result, source[rand.Intn(len(source))])
	}

	return string(result)
}
