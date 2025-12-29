package utils

import "golang.org/x/crypto/bcrypt"

func HashPass(password string) (string,error) {
	hash,err := bcrypt.GenerateFromPassword([]byte(password),14)

	return string(hash),err
}