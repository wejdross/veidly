package helpers

import (
	"fmt"
	"net/mail"
	"strings"
)

func ValidateEmail(email string) error {
	_, err := mail.ParseAddress(email)
	return err
}

func CRNG_EmailOrPanic() string {
	return strings.ToLower(
		fmt.Sprintf(CRNG_stringPanic(12) +
			"@" +
			CRNG_stringPanic(5) +
			"." +
			CRNG_stringPanic(3)))
}
