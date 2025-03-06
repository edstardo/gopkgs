package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(b), err
}
