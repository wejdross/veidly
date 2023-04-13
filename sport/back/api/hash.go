package api

import (
	"golang.org/x/crypto/bcrypt"
)

/*
	we may want to implement our custom hashing mechanism later
*/

/*
	Compare plaintext password and hash
*/
func CmpHash(pass string, hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(pass))
}

/*
	Generate password hash
*/
func GetHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}
