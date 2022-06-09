package helpers

import (
	bcrypt "golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func VerifyHashPassword(hashedPassword, password []byte) error {
	err := bcrypt.CompareHashAndPassword(hashedPassword, password)
	return err
}