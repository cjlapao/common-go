package helper

import "log"

// CheckError Checks if an error exist and calls panic
func CheckError(e error) {
	if e != nil {
		panic(e)
	}
}

// CheckErrorAndStop Checks if an error exist and stops the app
func CheckErrorAndStop(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// CheckEmptyStringAndStop Checks if an error exist and stops the app
func CheckEmptyStringAndStop(e string) {
	if len(e) != 0 {
		log.Fatal(e)
	}
}
