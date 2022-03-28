package util

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pwd string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash password fail %w", err)
	}
	return string(hashPassword), nil
}

func CheckPassword(pwd, hashPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(pwd))
}
