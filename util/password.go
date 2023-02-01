package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CreateHashedPwd(password string) (string, error) {
	hashedPwd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password.Error: %s", err)
	}
	return string(hashedPwd), nil
}

func ComparePwd(hashPwd string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashPwd), []byte(password))
}
